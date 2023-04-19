pragma solidity ^0.4.25;
pragma experimental ABIEncoderV2;

import "./Table.sol";
import "./Ownable.sol";
import "./MapStorage.sol";
import "./LibString.sol";
contract HistoryOrderInfo is Ownable {
    using LibString for string;
    MapStorage private mapStorage;
    TableFactory tf;
    string constant TABLE_NAME = "t_history_order_information3";
    constructor() public {
        tf = TableFactory(0x1001);
        tf.createTable(TABLE_NAME, "id","customerid,tradeYearMonth,financeId,data,key,hash,owner");
        mapStorage = new MapStorage();
    }
    
    event InsertOrderInfo(string customerid, string hash, string owner);
    event UpdateOrderInfo(string customerid, string hash, string owner);
    event HistoryOrderInfoHashNotFound(string customerid, string hash, string tips);
    
    function insert(string memory _id,string memory _params, string memory _data,string memory _key) public onlyOwner returns(int) {
        //params: tradeYearMonth financeId hash owner
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        string[] memory ss = _params.split(",");
        string memory _tradeYearMonth = ss[0];
        string memory _financeId = ss[1];
        string memory _hash = ss[2];
        string memory _owner = ss[3];
        require(!_isHashExist(table, _id,_hash), "current Hash has already exist");
        entry.set("customerid",_id);
        entry.set("tradeYearMonth",_tradeYearMonth);
        entry.set("financeId",_financeId);
        entry.set("data",_data);
        entry.set("key",_key);
        entry.set("hash",_hash);
        entry.set("owner",_owner);
        int256 count = table.insert(_id, entry);
        emit InsertOrderInfo(_id, _hash, _owner);
        return count;
    }
    function update(string memory _id,string memory _hash, string memory _owner) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        entry.set("owner",_owner);
        Condition condition = table.newCondition();
        condition.EQ("customerid", _id);
        condition.EQ("hash", _hash);
        int256 count = table.update(_id, entry, condition);
        if (count == 1){
            emit UpdateOrderInfo(_id, _hash, _owner);
        } else {
            emit HistoryOrderInfoHashNotFound(_id, _hash, "未找到记录");
        }
        return count;
    }
    function _isProcessIdExist(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        return _table.select(_id, condition).size() != int(0);
    }
    function _isHashExist(Table _table, string memory _id, string memory _hash) internal view returns(bool) {
        Condition condition = _table.newCondition();
        condition.EQ("hash", _hash);
        return _table.select(_id, condition).size() != int(0);
    }
    function select(string memory _id) private view returns(Entries _entries){
        Table table = tf.openTable(TABLE_NAME);
        require(_isProcessIdExist(table, _id), "HistoryOrderInfo select: current Id not exist");
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
            _json = _json.concat(_entry.getString("tradeYearMonth"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("financeId"));
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