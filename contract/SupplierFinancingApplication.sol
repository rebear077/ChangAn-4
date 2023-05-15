pragma solidity ^0.4.25;
pragma experimental ABIEncoderV2;

import "./Table.sol";
import "./Ownable.sol";
import "./MapStorage.sol";
import "./LibString.sol";
contract SupplierFinancingApplication is Ownable {
    MapStorage private mapStorage;
    using LibString for string;
    TableFactory tf;
    string constant TABLE_NAME = "t_supplier_financing_application";
    // 表名称：t_supplier_financing_application
    // 表主键：id 
    // 表字段：data
    // 字段含义：
    constructor() public {
        tf = TableFactory(0x1001);
        tf.createTable(TABLE_NAME, "id","financingid,customerid,data,key,hash,state");
        mapStorage = new MapStorage();
    }
    
    event InsertSupplierFinancingApplication(string id, string hash);
    event UpdateSupplierFinancingApplication(string id, string hash);
    event FinancingApplicationHashNotFound(string customerid, string hash, string tips);
    //发起融资申请
    function insert(string memory _id,string memory _customerID,string memory _data,string memory _key,string memory _hash) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        require(!_isFinanceidExist(table, _id), "current financingid has already exist");
        string memory _state="待审批";
        entry.set("financingid",_id);
        entry.set("customerid",_customerID);
        entry.set("data",_data);
        entry.set("key",_key);
        entry.set("hash",_hash);
        entry.set("state",_state);
        int256 count = table.insert(_id, entry);
        emit InsertSupplierFinancingApplication(_id, _hash);
        return count;
    }
    //更新融资申请
    function update(string memory _id, string memory _customerID, string memory _data,string memory _key,string memory _hash) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        require(_isFinanceidExist(table, _id), "SupplierFinancingApplication select: current financingId not exist");
        require(checkState(table,_id),"当前融资申请未被驳回");
        string memory _state="待审批";
        entry.set("financingid",_id);
        entry.set("customerid",_customerID);
        entry.set("data",_data);
        entry.set("key",_key);
        entry.set("hash",_hash);
        entry.set("state", _state);
        Condition condition = table.newCondition();
        int256 count = table.update(_id, entry,condition);
        if (count == 1){
            emit UpdateSupplierFinancingApplication(_id, _hash);
        } else {
            emit FinancingApplicationHashNotFound(_id, _hash, "未找到记录");
        }
        return count;
    }
    function _isFinanceidExist(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        condition.EQ("id",_id); 
        condition.EQ("financingid",_id);
        return _table.select(_id, condition).size() != int(0);
    }
    function select(string memory _id) private view returns(Entries _entries){
        Table table = tf.openTable(TABLE_NAME);
        require(_isFinanceidExist(table, _id), "SupplierFinancingApplication select: current processId not exist");
        Condition condition = table.newCondition();
        _entries = table.select(_id, condition);
        return _entries;
    }
    function checkState(Table _table,string memory _id)internal view returns(bool){
        Condition condition = _table.newCondition();
        Entries _entries=_table.select(_id,condition);
        Entry _entry=_entries.get(0);
        string memory result=_entry.getString("state");
        string memory flag="驳回";
        return LibString.equal(result,flag);
    }
    function getDetail(string memory _id) public view returns(string memory _json){
        Entries _entries = select(_id);
        _json = _returnData(_entries);
    }
    function _returnData(Entries _entries) internal view returns(string){

        string memory _json = "{";
        for (int256 i=0;i<_entries.size();i++){
            Entry _entry=_entries.get(i);
            _json=_json.concat("[");
            _json = _json.concat(_entry.getString("financingid"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("customerid"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("data"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("key"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("hash"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("state"));
            _json = _json.concat("]");
        }
        _json=_json.concat("}");
        return _json;
    }
}