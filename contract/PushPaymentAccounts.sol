pragma solidity ^0.4.25;
pragma experimental ABIEncoderV2;

import "./Table.sol";
import "./Ownable.sol";
import "./MapStorage.sol";
import "./LibString.sol";
contract PushPaymentAccounts is Ownable {
    using LibString for string;
    MapStorage private mapStorage;
    TableFactory tf;
    string constant TABLE_NAME = "t_push_payment_accounts1";
    // 表名称：t_push_payment_accounts
    // 表主键：id 
    // 表字段：data
    // 字段含义：
    constructor() public {
        tf = TableFactory(0x1001);
        tf.createTable(TABLE_NAME, "id","financeID,data,key,hash,state");
        mapStorage = new MapStorage();
    }
    
    event UpdatePaymentAccounts(string customerid, string hash);
    event PaymentAccountsHashNotFound(string customerid, string hash, string tips);
    function lock(string memory _id,string memory _financeID,string memory _hash)public onlyOwner returns(int){
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        
        string memory _state="未锁定";
        string memory _newstate="已锁定";
        require(check(table,_id,_financeID,_hash),"check failed");
        Condition condition = table.newCondition();
        condition.EQ("hash",_hash);
        condition.EQ("financeID",_financeID);
        condition.EQ("state",_state);
        entry.set("state",_newstate);
        int256 count = table.update(_id, entry,condition);
        if (count == 1){
            emit UpdatePaymentAccounts(_id, _hash);
        } else {
            emit PaymentAccountsHashNotFound(_id, _hash, "未找到记录");
        }
        return count;
    }
    function updateAndLock(string memory _idAndFinanceID,string memory _data,string memory _key,string memory _newhash,string memory _oldhash) public onlyOwner returns(int) {
        string[] memory ss = _idAndFinanceID.split(",");
        // string memory _id = ss[0];
        // string memory _financeID = ss[1];
        Table table = tf.openTable(TABLE_NAME);
        Entry entry = table.newEntry();
        require(check(table,ss[0],ss[1],_oldhash),"check failed");
        string memory _state="未锁定";
        string memory _newstate="已锁定";
        Condition condition = table.newCondition();
        condition.EQ("hash",_oldhash);
        condition.EQ("financeID",ss[1]);
        condition.EQ("state",_state);
        entry.set("data", _data);
        entry.set("key",_key);
        entry.set("hash",_newhash);
        entry.set("state",_newstate);
        int256 count = table.update(ss[0], entry,condition);
        if (count == 1){
            emit UpdatePaymentAccounts(ss[0], _newhash);
        } else {
            emit PaymentAccountsHashNotFound(ss[0], _oldhash, "未找到记录");
        }
        return count;
    }
    function check(Table _table,string memory _id,string memory _financeID,string memory _oldhash)internal view returns(bool){
        require(_isProcessIdExist(_table, _id), "current customerid not exist");
        require(isHashExisting(_table,_id,_oldhash),"old hash not exist");
        string memory _state="未锁定";
        Condition condition = _table.newCondition();
        condition.EQ("hash",_oldhash);
        condition.EQ("financeID",_financeID);
        condition.EQ("state",_state);
        _table.select(_id,condition);
        Entries _entries=_table.select(_id,condition);
        return _entries.size()!=int(0);
    }
    function _isProcessIdExist(Table _table, string memory _id) internal view returns(bool) {
        Condition condition = _table.newCondition();
        return _table.select(_id, condition).size() != int(0);
    }
    function select(string memory _id) private view returns(Entries _entries){
        Table table = tf.openTable(TABLE_NAME);
        require(_isProcessIdExist(table, _id), "PushPaymentAccounts select: current processId not exist");
        Condition condition = table.newCondition();
        _entries = table.select(_id, condition);
        return _entries;
    }
    function getDetail(string memory _id) public view returns(string memory _json){
        Entries _entries = select(_id);
        _json = _returnData(_entries);
    }
   
    function isHashExisting(Table _table,string memory _id, string memory _hash)internal view returns(bool){
        Entry entry = _table.newEntry();
        Condition condition = _table.newCondition();
        condition.EQ("hash",_hash);
        Entries _entries=_table.select(_id,condition);
        return _entries.size()!=int(0);
    }
    function _returnData(Entries _entries) internal view returns(string){

        string memory _json = "{";
        for (int i=0;i<_entries.size();i++){
            Entry _entry=_entries.get(i);
            _json=_json.concat("[");
            _json = _json.concat(_entry.getString("id"));
            _json = _json.concat(",");
            _json = _json.concat(_entry.getString("financeID"));
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