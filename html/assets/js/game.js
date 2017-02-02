/**
 * 游戏控制
 */
;(function(){
    /***************************************/
    /* view
    /***************************************/

    var loadLoginForm = function(){
        console.log('load login form');
        U.tpl('login', function($root){
            // click to login
            $root.on('click', 'input[name=login]', function(){
                var $username = $root.find('input[name=username]');
                var $password = $root.find('input[name=password]');
                login($username.val(), $password.val());
            });
            // login(username, password)
            $root.on('click', 'input[name=register]', function(){
                loadRegisterForm();
            });
        });
    };
    loadRegisterForm = function(){
        console.log('load register form');
        U.tpl('register', function($root){
            $root.on('click', 'input[name=return]', function(){
                loadLoginForm();
            });
        });
    };

    var loadGameList = function(data){
        console.log('load game list' , JSON.stringify(data));
        U.tpl('desk', data, function($root){
            // show not on desk
            // click to creat or join into
            $root.on('click', 'input[name=join]', function(){
                var fsmId = "" + $(this).data('id');
                if (fsmId) {
                    join(fsmId);
                }
            });
            $root.on('click', 'input[name=create]', function(){
                create();
            });
        });
    };
    var loadReadyView = function(data){
        console.log('load ready view' , JSON.stringify(data));
        U.tpl('ready', data, function($root){
            // waiting for user ready
            $root.on('click', 'input[name=leave]', function(){
                leave();
            });
        });
    };
    var loadNextView = function(data){
        console.log('load next view' , JSON.stringify(data));
        data.selfPlayer = data.gameInfo.playerList[C.user_id];
        U.tpl('result', data, function($root){
            // waiting for user click next
            $root.on('click', 'input[name=next]', function(){
                next();
            });
        });
    };

    var loadPlayingView = function(data, with_click){
        console.log('load playing view' , JSON.stringify(data));
        if (data.gameInfo.maxPlayer == 2) {
            _.map(data.gameInfo.playerChain, function(playerId, seq){
                if (playerId == C.user_id) {
                    data.selfPlayer = data.gameInfo.playerList[playerId];
                    data.facePlayer = data.gameInfo.playerList[data.gameInfo.playerChain[(seq + 1)%2]];
                    data.lastPlayer = null;
                    data.nextPlayer = null;
                }
            });
        } else if (data.gameInfo.maxPlayer == 4) {
            _.map(data.gameInfo.playerChain, function(playerId, seq){
                if (playerId == C.user_id) {
                    data.selfPlayer = data.gameInfo.playerList[playerId];
                    data.nextPlayer = data.gameInfo.playerList[data.gameInfo.playerChain[(seq + 1)%4]];
                    data.facePlayer = data.gameInfo.playerList[data.gameInfo.playerChain[(seq + 2)%4]];
                    data.lastPlayer = data.gameInfo.playerList[data.gameInfo.playerChain[(seq + 3)%4]];
                }
            });
        }
        U.tpl('playing', data, function($root){
            // update info to show
            // show op button
            var multi_sel = 0;
            var sel_cache = [];
            var cmd_cache = null;
            if (!with_click) {
                return;
            }
            // waiting for user op
            // click to request op with cmd
            // op(cmd, sel)
            $root.on('click', '.player-0 input.card', function(){
                var index = $(this).data('index');
                if (!multi_sel) {
                    var sel = [];
                    if (index || index == 0) {
                        sel.push(index);
                    }
                    op(7, sel);
                } else {
                    if (sel_cache.length == multi_sel) {
                        return;
                    } else if (sel_cache.length < multi_sel) {
                        sel_cache.push(index);
                        if (sel_cache.length == multi_sel) {
                            op(cmd_cache, sel_cache);
                        }
                    }
                }
            });
            $root.on('click', 'input.op', function(){
                var cmd = $(this).data('op');
                // TODO switch cmd commit sel_cache
                switch (cmd) {
                    case 12:
                    case 2:
                    case 5:
                    case 3:
                        op(cmd);
                        return;
                    case 6:
                        cmd_cache = cmd;
                        multi_sel = 2;
                        return;
                }
            });
            $root.on('click', 'input[name=leave]', function(){
                leave();
            });
        });
    };

    /***************************************/
    /* controller
    /***************************************/
    var game_info = {};
    var login = function(username, password){
        var user_list = {
            "yy":"1",
            "jj":"2",
            "mm":"3",
            "bb":"4"
        };
        if (user_list[username] && password == "123") {
            C.user_id = user_list[username];
            checkGameStatus();
        }
    };

    var register = function(){};

    var checkUserStatus = function(){
        if (C.user_id) {
            checkGameStatus();
        } else {
            loadLoginForm();
        }
    };

    var checkGameStatusPid = 0;
    var checkGameStatus = function(){
        U.action('info', {
            success: function(res){
                if (!res.gameInfo) {
                    loadGameList(res);
                    setTimeout(function(){
                        checkGameStatus();
                    }, 1000);
                } else if (res.gameInfo.setStatus == 1) {
                    loadReadyView(res);
                    setTimeout(function(){
                        checkGameStatus();
                    }, 1000);
                } else if (res.gameInfo.setStatus == 3) {
                    loadNextView(res);
                    setTimeout(function(){
                        checkGameStatus();
                    }, 1000);
                } else if (res.gameInfo.curEvent.PlayerId == C.user_id) {
                    loadPlayingView(res, true);
                } else {
                    loadPlayingView(res, false);
                    setTimeout(function(){
                        checkGameStatus();
                    }, 500);
                }
            }
        });
    };

    var next = function(){
        U.action('next', {
            success:function(){
                //checkGameStatus();
            }
        });
    };

    var leave= function(){
        if (confirm("放弃这场比赛?")) {
            U.action('leave', {
                success:function(){
                    //checkGameStatus();
                }
            });
        }
    };

    var create = function(fsmId){
        U.action('create', {
            success:function(res){
                // waiting for status refresh
                //checkGameStatus();
            }
        });
    };

    var join = function(fsmId) {
        U.action('join', {
            data:{
                fsmId:fsmId
            },
            success:function(){
                checkGameStatus();
            }
        });
    }

    var op = function(cmd, sel){
        U.action('op', {
            data: {
                op: cmd,
                indexs: sel||[]
            },
            success:function(){
                // waiting for status refresh
                checkGameStatus();
            }
        });
    };

    /***************************************/
    /* init
    /***************************************/
    var init = function(){
        checkUserStatus();
    };
    init();

})();
