//复制链接
$('#copy').off('click').on('click', function(){
    var url = $('.url').val();
    if(!url){
        alert('您复制的链接为空！');
        return;
    }
   $('.url').select();
   document.execCommand('Copy');
   alert("复制成功");
});