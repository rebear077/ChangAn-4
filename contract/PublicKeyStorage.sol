pragma solidity ^0.4.25;
pragma experimental ABIEncoderV2;
import "./LibString.sol";
import "./Table.sol";
import "./Ownable.sol";
import "./MapStorage.sol";
contract PublicKeyStorage is Ownable {
    using LibString for string;
    MapStorage private mapStorage;
    TableFactory tf;
    string constant TABLE_NAME = "t_public_key";
    // 表名称：t_public_key
    // 表主键：id 
    // 表字段：role,key
    // 字段含义：
    constructor() public {
        tf = TableFactory(0x1001);
        tf.createTable(TABLE_NAME, "id","role,key");
        mapStorage = new MapStorage();
    }
    function insert(string memory _id, string memory _role,string memory _key) public onlyOwner returns(int) {
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        entry.set("role",_role);
        entry.set("key",_key);
        int256 count = table.insert(_id, entry);
        return count;
    }
    function _isProcessIdExist(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        return _table.select(_id, condition).size() != int(0);
    }
    function select(string memory _id) private view returns(Entries _entries){
        Table table = tf.openTable(TABLE_NAME);
        require(_isProcessIdExist(table, _id), "PublicKeyStorage select: current processId not exist");
        Condition condition = table.newCondition();
        _entries = table.select(_id, condition);
        return _entries;
    }
    function getDetail(string memory _id) public view returns(string memory _json){
        Entries _entries = select(_id);
        _json = _returnData(_entries);
    }
    function _returnData(Entries _entries) internal view returns(string){
        string memory _json;
        for (int i=0;i<_entries.size();i++){
            Entry _entry=_entries.get(i);
            _json = _json.concat(_entry.getString("key"));
         
        }
        return _json;
    }

}