
//分页
$('#pagination').my_page('#searchForm');


//添加账号
$('.add').off('click').on('click', function(){
    $('.password').show();
    $('#c-userId').val('');  //账号ID
    $('#c-account').val('').removeAttr('disabled');  //账号
    $('#c-password').val('');  //密码
    $('#c-user').val('');     //使用人
    $('#c-post').val(0);     //角色
    $('#c-accountState').val(2);     //账号状态
    layer.open({
        type: 1,
        title: '添加账号',
        area: '500px',
        content: $('#content'),
        cancel: function(){
            formError.hide();
        }
    });
});

//取消按钮
$('#cancel').off('click').on('click', function(){
    formError.hide();
    layer.closeAll();
});

//表单的错误提示信息
var formError = {
    show: function(errMsg){
        $('.error-item').show().find('.error-msg').text(errMsg);
    },
    hide: function(){
        $('.error-item').hide().find('.error-msg').text('');
    }
}

//点击确定
$('#save').off('click').on('click', function(){
    var c_userId = parseInt($('#c-userId').val());
    var c_account = $.trim($('#c-account').val());
    var c_password = $.trim($('#c-password').val());
    var c_user = $.trim($('#c-user').val());
    var c_post = parseInt($('#c-post').val());
    var c_accountState = parseInt($('#c-accountState').val());
    if(!c_account){
        formError.show('账号不得为空！');
        return;
    }
    if(!c_userId){    //添加账号
        if(!c_password){
            formError.show('密码不得为空！');
            return;
        }
        if(c_password && c_password.length < 6){
            formError.show('密码不得少于6位！');
            return;
        }
    }
    
    if(!c_user){
        formError.show('使用人姓名不得为空！');
        return;
    }
    if(c_post == 0){
        formError.show('请选择角色！');
        return;
    }
    if(c_accountState == 2){
        formError.show('请选择账号状态！');
        return;
    }
    //验证全部通过
    $.postJSON('../admin/useredit',
    {
        Id: c_userId,
        Account: c_account,
        DisplayName: c_user,
        RoleId: c_post,
        Password: c_password,
        IsUsed: c_accountState
    },
    function(res){
        if (res.ret==200){
            alert("SUCCESS!");
            layer.closeAll();
            getpage("../admin/userlist");
        }else{
            alert(res.err);
        };
        }, 'post');

});


//编辑账号
$('.edit').off('click').on('click', function(){
    $('.password').hide();
    $('#c-userId').val($(this).parent().siblings('.userId').text());  //账号ID
    $('#c-account').val($(this).parent().siblings('.account').text()).attr('disabled', true);  //账号
    $('#c-user').val($(this).parent().siblings('.user').text());     //使用人
    $('#c-post').val($(this).parent().siblings('.post').attr('data-index'));     //角色
    $('#c-accountState').val($(this).parent().siblings('.accountState').attr('data-index'));   //账号状态
    layer.open({
        type: 1,
        title: '编辑账号',
        area: '500px',
        content: $('#content'),
        cancel: function(){
            formError.hide();
        }
    });
});



//删除账号
$('.delete').off('click').on('click', function(){
    var userId = parseInt($(this).parent().siblings('.userId').text());
    layer.confirm('你确定要删除该账号吗？', {
        btn: ['确定', '取消']
    }, function(){
        $.postJSON('../admin/deluser', {
            Id: userId
        },
        function(res){
            if (res.ret==200){
                alert("删除成功！");
                layer.closeAll();
                getpage("../admin/userlist");
            };
        }, 'post');
    });
});
