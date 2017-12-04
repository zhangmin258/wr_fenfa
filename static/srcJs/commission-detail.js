



//分页
$("#pagination").my_page("#searchForm");

//提示文字说明
$('.ps-1').on('click', function(){
    layer.open({
        type: 1,
        shade: false,
        title: false, 
        content: $('.layer_notice1'),
    });
});


//申请提现
$('#withdrawal-application').on('click', function(){
    layer.open({
        type: 1,
        title: '申请提现',
        area: '500px',
        content: $('#content'),
        cancel: function(){
            $('#withdrawal-amount').val('');
        }
    });
});

var maxCount = $('.data4').text();
$('#withdrawal-amount').attr('placeholder', '当前可提现金额为'+maxCount+'元');

//取消按钮
$('#cancel').on('click', function(){
    layer.closeAll();
    $('#withdrawal-amount').val('');
});



//点击保存按钮
$('#save').off('click').on('click', function(){
    var bankCardId = parseInt($('#bankCardId').text());
    var amount = parseFloat($('#withdrawal-amount').val());
    var maxCount = parseFloat($('.data4').text());
    if(!bankCardId){
        alert('银行卡信息出错！');
        return;
    }
    if (!amount ){
        alert('请输入提现金额！');
        return;
    };
    if(amount && amount < 100){
        alert('最小提现金额不得少于100元，当前不能提现！');
        return;
    }
    if (amount && amount > maxCount){
        alert('您已超出最大提现金额，无法提现！');
        return;
    }
    $.postJSON('../withdraw/withdrawcash',
        {
            BankCardId: bankCardId,
            Money: amount
        },
        function(res){
            if (res.ret == 200){
                layer.closeAll();
                alert("提交成功，等待人工审核！");
                getpage('../commission/commissioninfo');
            }else{
                alert(res.err);
            };
        }, 'post');
});


//点击搜索
$('#submit').off('click').on('click', function(){
    $('#searchForm').submit();
});


$('#withdrawal-amount').off().on('keyup', function(){
    var amount = parseFloat($(this).val()); 
    if(amount < 0){
        $('#withdrawal-amount').val(0);
        return;
    }
});

