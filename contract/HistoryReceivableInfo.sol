pragma solidity ^0.4.25;
pragma experimental ABIEncoderV2;

import "./Table.sol";
import "./Ownable.sol";
import "./MapStorage.sol";
import "./LibString.sol";
contract HistoryReceivableInfo is Ownable {
    using LibString for string;
    MapStorage private mapStorage;
    TableFactory tf;
    string constant TABLE_NAME = "t_history_receivable_information";
    constructor() public {
        tf = TableFactory(0x1001);
        tf.createTable(TABLE_NAME, "id","time,data,key,hash");
        mapStorage = new MapStorage();
    }
    function insert(string memory _id,string memory _time, string memory _data,string memory _key,string memory _hash) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        entry.set("time",_time);
        entry.set("data",_data);
        entry.set("key",_key);
        entry.set("hash",_hash);
        int256 count = table.insert(_id, entry);
        return count;
    }
     function _isProcessIdExist(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        return _table.select(_id, condition).size() != int(0);
    }
    function select(string memory _id) private view returns(Entries _entries){
        Table table = tf.openTable(TABLE_NAME);
        require(_isProcessIdExist(table, _id), "HistoryUsedInfo select: current processId not exist");
        Condition condition = table.newCondition();
        _entries = table.select(_id, condition);
        return _entries;
    }
    function getDetailInList(string memory _id) public view returns(string memory _json){
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
            _json = _json.concat(_entry.getString("time"));
            _json = _json.concat(",");
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
            _json=_json.concat("\"time\":\"");
            _json = _json.concat(_entry.getString("time"));
            _json = _json.concat("\",");
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