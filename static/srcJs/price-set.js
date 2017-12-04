

//分页
$('#pagination').my_page('#searchForm');


//设置价格
$('.edit').off('click').on('click', function(){
    var productId = $(this).parent().siblings('.productId').text();   //产品Id
    var productName = $(this).parent().siblings('.productName').text();   //产品名称
    var productPrice = $(this).parent().siblings('.productPrice').text();   //产品价格
    $('#productId').val(productId);
    $('#productName').val(productName).attr('disabled', true);
    $('#price').val(productPrice);
    layer.open({
        type: 1,
        title: '设置价格',
        area: '500px',
        content: $('#content')
    });
});


//点击取消按钮
$('#cancel').off('click').on('click', function(){
    layer.closeAll();
});


//点击保存按钮
$('#save').off('click').on('click', function(){
    var productId = parseInt($('#productId').val());
    var productPrice = parseFloat($('#price').val());
    if(!productPrice){
       alert('价格不得为空！');
       return;
    }
    if(productPrice && productPrice < 0){
        alert('价格不得小于零！');
        return;
    }
    $.postJSON('../admin/saveproduct',
        {
            ProductId: productId,
            AgentPrice: productPrice
        },
        function(res){
            if (res.ret==200){
                layer.closeAll();
                alert("修改成功！");
                getpage('../admin/pricelist');
            }else{
                alert(res.err);
            }
        }, 'post');
});
