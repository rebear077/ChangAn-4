pragma solidity ^0.4.25;
pragma experimental ABIEncoderV2;

import "./Table.sol";
import "./Ownable.sol";
import "./MapStorage.sol";
import "./LibString.sol";
contract InvoiceInformationStorage is Ownable {
    using LibString for string;
    MapStorage private mapStorage;
    TableFactory tf;
    string constant TABLE_NAME = "t_invoice_information3";
    // 表名称：t_invoice_information
    // 表主键：id 
    // 表字段：data
    // 字段含义：
    constructor() public {
        tf = TableFactory(0x1001);
        tf.createTable(TABLE_NAME, "id","customerid,time,type,num,data,key,hash,owner");
        mapStorage = new MapStorage();
    }
    
    event InsertInvoiceInfo(string customerid, string hash, string owner);
    event UpdateInvoiceInfo(string customerid, string hash, string owner);
    event InvoiceInfoHashNotFound(string customerid, string hash, string tips);
    
    function insert(string memory _id, string memory _params, string memory _data,string memory _key) public onlyOwner returns(int) {
        string[] memory ss = _params.split(",");
        string memory _time = ss[0];
        string memory _type = ss[1];
        string memory _num = ss[2];
        string memory _hash = ss[3];
        string memory _owner = ss[4];
        Table table = tf.openTable(TABLE_NAME);
        require(!_isHashExist(table, _id, _hash), "current Hash has already exist");
        Entry entry = table.newEntry();
        entry.set("customerid",_id);
        entry.set("time",_time);
        entry.set("type",_type);
        entry.set("num",_num);
        entry.set("data",_data);
        entry.set("key",_key);
        entry.set("hash",_hash);
        entry.set("owner",_owner);
        int256 count = table.insert(_id, entry);
        emit InsertInvoiceInfo(_id, _hash, _owner);
        return count;
    }
    function update(string memory _id, string memory _hash,string memory _owner) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        require(_isOwnerNull(table, _id), "The invoice has been used");
        entry.set("owner",_owner);
        Condition condition = table.newCondition();
        require(_isHashExist(table, _id, _hash), "Hash not exists");
        condition.EQ("customerid", _id);
        condition.EQ("hash", _hash);
        int256 count = table.update(_id, entry,condition);
        if (count == 1) {
            emit UpdateInvoiceInfo(_id, _hash, _owner);
        } else {
            emit InvoiceInfoHashNotFound(_id, _hash, "未找到记录");
        }
        
        return count;
    }
    function _isProcessIdExist(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        return _table.select(_id, condition).size() != int(0);
    }
    function _isHashExist(Table _table, string memory _id, string memory _hash) internal view returns(bool) {
        Condition condition = _table.newCondition();
        condition.EQ("hash",_hash);
        return _table.select(_id, condition).size() != int(0);
    }
    function _isOwnerNull(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        condition.EQ("owner","");
        return _table.select(_id, condition).size() != int(0);
    }
    function select(string memory _id) private view returns(Entries _entries){
        Table table = tf.openTable(TABLE_NAME);
        require(_isProcessIdExist(table, _id), "InvoiceInformationStorage select: current processId not exist");
        Condition condition = table.newCondition();
        _entries = table.select(_id, condition);
        return _entries;
    }
    function getDetailInList(string memory _id) public view returns(string memory _json){
        Entries _entries = select(_id);
        _json = _returnData(_entries);
    }
    function _returnData(Entries _entries) internal view returns(string){
        string memory _json = "{";
        for (int256 i=0;i<_entries.size();i++){
            Entry _entry=_entries.get(i);
            _json=_json.concat("[");
            _json = _json.concat(_entry.getString("customerid"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("time"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("type"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("num"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("data"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("key"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("hash"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("owner"));
            _json = _json.concat("]");
        }
        _json=_json.concat("}");
        return _json;
    }
}