
var rootSetting = {
    view: {
        showIcon: false
    },
    check: {
        enable: true
    },
    data: {
        simpleData: {
            enable: true,
            rootPId: 0
        }
    }
};

$.postJSON('/system/menudata', {
    RoleId: 'all'
}, function(result){
    if (!!result && result.ret == 200) {
        $.fn.zTree.init($("#treeRoot"), rootSetting, result.m);
        $.postJSON('/system/menudata', {
            RoleId:'visit'
        }, function(result){
            if (!!result && result.ret == 200) {
                var treeMenuObj = $.fn.zTree.getZTreeObj("treeRoot");
                treeMenuObj.checkAllNodes(false);
                if (!!result.m) {
                    for (var i = 0; i < result.m.length; i++) {
                        var nodes = treeMenuObj.getNodesByParam("id", result.m[i].id, null);
                        if (!!nodes && nodes.length) {
                            treeMenuObj.checkNode(nodes[0], true, false);
                        }
                    };
                }
            }else {
                alert('没有数据.')
            }
        }, 'post');

    }else if(!!result.msg) {
        alert(result.msg);
    }else{
        alert('没有数据');
    }
}, 'post');

function Post_form_back(){
    return getpage('/system/rolelist');
}

function Post_form_submit(){
    var account = $("#account").val();
    var rid = $("#rid").val();
    $(".error").html("");
    if(!account){
        $(".account_num_error").html("请填写账号");
    }else{
        var treeObj = $.fn.zTree.getZTreeObj("treeRoot");
        var nodes = treeObj.getCheckedNodes(true);
        if(nodes.length>0){
            var checkIdArr = [];
            for(var i=0; i<nodes.length; i++){
                checkIdArr.push(nodes[i].id);
            }
            $.postJSON('/system/visitmenuedit', {
                Rid:rid,
                Account:account,
                CheckId:checkIdArr.join(",")
            }, function(result){
                if(!!result && result.ret == 200){ // 成功
                    alert('成功');
                    Post_form_back();
                }else if(!!result){ // 失败
                    $(".account_num_error").html(result.err);
                } else {
                    alert('没有数据返回');
                }
            }, 'post');
        }else{
            alert("请勾选权限")
            $(".root_error").html("请勾选权限");
        }
    }    
}