
//分页
$('#pagination').my_page('#searchForm');

//点击搜索
$('#submit').off('click').on('click', function(){
    // var phone = $.trim($("#phone").val());
    // if (phone && !(/^1[34578]\d{9}$/.test(phone))) {
    //     alert('手机号码格式不正确！');
    //     return;
    // }
    $('#searchForm').submit();
});


//通过提现申请
$('.enable').on('click', function(){
    var depId = $(this).parent().siblings('.depId').text();
    layer.confirm('确定通过提现申请？', {
        btn: ['确定', '取消']
    }, function(){
        $.postJSON('../withdraw/withdrawdeposit',
            {
                DepId: depId
            },
            function(res){
                if (res.ret==200){
                    layer.closeAll();
                    alert("操作成功！");
                    getpage('../withdraw/depositlist');
                }else{
                    alert(res.err);
                };
            }, 'post');
        layer.closeAll();
    })
});


//拒绝提现申请
$('.disable').on('click', function(){
    var depId = $(this).parent().siblings('.depId').text();
    layer.confirm('确定拒绝提现申请？', {
        btn: ['确定', '取消']
    }, function(){
        $.postJSON('../withdraw/refusewithdrawdeposit',
            {
                DepId: depId
            },
            function(res){
                if (res.ret==200){
                    layer.closeAll();
                    alert("操作成功！");
                    getpage('../withdraw/depositlist');
                }else{
                    alert(res.err);
                };
            }, 'post');
        layer.closeAll();
    })
});

//点击取消按钮
$('#cancel').off('click').on('click', function(){
    layer.closeAll();
});

if(window.sessionStorage.phone){
    $('#phone').val(window.sessionStorage.phone);
}