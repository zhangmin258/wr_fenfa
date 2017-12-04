

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


if(window.sessionStorage.phone){
    $('#phone').val(window.sessionStorage.phone);
}


