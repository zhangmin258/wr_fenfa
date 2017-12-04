
//文字说明
$('.ps').on('click', function(){
    layer.open({
        type: 1,
        shade: false,
        title: false, 
        content: $('.layer_notice')
    });
});

//分页
$("#pagination").my_page("#searchForm");