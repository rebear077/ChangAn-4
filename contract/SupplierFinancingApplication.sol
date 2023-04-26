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
    string constant TABLE_NAME = "t_supplier_financing_application1";
    // 表名称：t_supplier_financing_application
    // 表主键：id 
    // 表字段：data
    // 字段含义：
    constructor() public {
        tf = TableFactory(0x1001);
        tf.createTable(TABLE_NAME, "id","financingid,data,key,hash");
        mapStorage = new MapStorage();
    }
    
    event InsertSupplierFinancingApplication(string id, string financingid, string hash);
    event UpdateSupplierFinancingApplication(string id, string financingid, string hash);
    event FinancingApplicationHashNotFound(string customerid, string hash, string tips);
    
    function insert(string memory _id, string memory _financingid, string memory _data,string memory _key,string memory _hash) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        require(_isFinanceidExist(table, _id, _hash), "current financingid or Hash has already exist");
        entry.set("financingid",_financingid);
        entry.set("data",_data);
        entry.set("key",_key);
        entry.set("hash",_hash);
        int256 count = table.insert(_id, entry);
        emit InsertSupplierFinancingApplication(_id, _financingid, _hash);
        return count;
    }
    function update(string memory _id, string memory _financingid, string memory _data,string memory _key,string memory _hash) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        require(_isProcessIdExist(table, _id), "SupplierFinancingApplication select: current financingId not exist");
        entry.set("financingid",_financingid);
        entry.set("data",_data);
        entry.set("key",_key);
        entry.set("hash",_hash);
        int256 count = table.insert(_id, entry);
        if (count == 1){
            emit UpdateSupplierFinancingApplication(_id, _financingid, _hash);
        } else {
            emit FinancingApplicationHashNotFound(_id, _hash, "未找到记录");
        }
        return count;
    }
    function _isProcessIdExist(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        return _table.select(_id, condition).size() != int(0);
    }
    function _isFinanceidExist(Table _table, string memory _id, string memory _hash) internal view returns(bool) {
        Condition condition = _table.newCondition();
        condition.EQ("id",_id);
        condition.EQ("financingid",_id);
        condition.EQ("hash",_hash);
        return _table.select(_id, condition).size() == int(0);
    }
    function select(string memory _id) private view returns(Entries _entries){
        Table table = tf.openTable(TABLE_NAME);
        require(_isProcessIdExist(table, _id), "SupplierFinancingApplication select: current processId not exist");
        Condition condition = table.newCondition();
        _entries = table.select(_id, condition);
        return _entries;
    }
    function getDetail(string memory _id) public view returns(string memory _json){
        Entries _entries = select(_id);
        _json = _returnData(_entries);
    }
    function getDetailInJson(string memory _id) public view returns(string memory _json){
        Entries _entries = select(_id);
        _json = _returnJson(_entries);
    }
    function _returnData(Entries _entries) internal view returns(string){

        string memory _json = "{";
        for (int256 i=0;i<_entries.size();i++){
            Entry _entry=_entries.get(i);
            _json=_json.concat("[");
            _json = _json.concat(_entry.getString("data"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("key"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("hash"));
            _json = _json.concat("]");
        }
        _json=_json.concat("}");
        return _json;
    }
    function _returnJson(Entries _entries)internal view returns(string){

        string memory _json = "[";
        for (int256 i=0;i<_entries.size();i++){
            Entry _entry=_entries.get(i);
            _json=_json.concat("{");
            _json=_json.concat("\"data\":\"");
            _json = _json.concat(_entry.getString("data"));
            _json = _json.concat("\",");
            _json=_json.concat("\"key\":\"");
            _json = _json.concat(_entry.getString("key"));
            _json = _json.concat("\",");
            _json=_json.concat("\"hash\":\"");
            _json = _json.concat(_entry.getString("hash"));
            _json = _json.concat("\"}");
            if (i!=_entries.size()-1){
              _json =_json.concat(",");  
            }
        }
        _json=_json.concat("]");
        return _json;
    }
    
}