
//分页
$('#pagination').my_page('#searchForm');

var userCode = $('#userCode').val();
$.each($('.code'), function(index, value){
    if($(value).text() === userCode){
        $(value).siblings('.operation').html('<span class="unomal">不可操作</span>');
        $(value).siblings('.userFc').html('-');
    }
});

//1为管理员，不能编辑
var AccountType = $('#AccountType').val();
if(AccountType == 1){
    $('.edit').hide();
}else{
    $('.edit').show();
}

//点击编辑
$('.edit').off('click').on('click', function(){
    var userId = $(this).parent().siblings('.userId').text();
    var code = $(this).parent().siblings('.code').text();
    var userName = $(this).parent().siblings('.userName').text();
    var phoneNum = $(this).parent().siblings('.phoneNum').text();
    var userFc = $(this).parent().siblings('.userFc').find('em').text();
    $('#userId').val(userId);
    $('#code').val(code);
    $('#user-name').val(userName);
    $('#phone-num').val(phoneNum).attr('disabled', true);
    $('#user-fc').val(userFc);
    layer.open({
        type: 1,
        title: '编辑员工信息',
        area: '500px',
        content: $('#content')
    });
});

//取消按钮
$('#cancel').off('click').on('click', function(){
    layer.closeAll();
});

//点击保存
$('#save').off('click').on('click', function () {
    var userName = $.trim($('#user-name').val());
    var userFc = $.trim($('#user-fc').val());
    var userAgentId = $.trim($('#userId').val());
    var code = $.trim($('#code').val());
    if(!userName){
        alert('姓名不得为空！');
        return;
    }
    if(!userFc){
        alert('员工分成不得为空！');
        return;
    }
    if(userFc && userFc < 0){
        alert('员工分成不得小于零！');
        return;
    }
    $.postJSON('../staff/savestaffinfo', {
        Code: code,
        Remark: userName,
        PriceScale: userFc,
        UserAgentId: userAgentId
    }, function(res){
        if(res.ret==200){
            alert("保存成功");
            layer.closeAll();
            getpage('/staff/stafflist');
        }else{
            alert(res.err);
        }
    }, 'post');
});



//点击禁用登录
$('.disable').off('click').on('click', function () {
    var code = $(this).parent().siblings('.code').text();
    layer.confirm('确定要禁用登录吗？', {
        btn: ['确定', '取消']
    }, function () {
        $.postJSON('../staff/locklogin', {
            Code: code,
            State: 1
        }, function(res){
            if (res.ret == 200) {
                alert("禁用成功");
                layer.closeAll();
                getpage('/staff/stafflist');
            } else {
                alert(res.err)
            }
        }, 'post');
    });
});


//点击启用登录
$('.enable').off('click').on('click', function () {
    var code = $(this).parent().siblings('.code').text();
    layer.confirm('确定要启用登录吗？', {
        btn: ['确定', '取消']
    }, function () {
        $.postJSON('../staff/locklogin', {
            Code: code,
            State: 0
        }, function(res){
            if (res.ret == 200) {
                alert("启用成功");
                layer.closeAll();
                getpage('/staff/stafflist');
            } else {
                alert(res.err)
            }
        }, 'post');
    });
});

//点击搜索
$('#submit').off('click').on('click', function(){
    // var phone = $.trim($("#phone").val());
    // if (phone && !(/^1[34578]\d{9}$/.test(phone))) {
    //     alert('手机号码格式不正确！');
    //     return;
    // }
    $('#searchForm').submit();
});
