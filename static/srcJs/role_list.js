

//删除岗位
$('.delete').off('click').on('click', function(){
    var postId = String($(this).data('id'));
    layer.confirm('删除该角色会造成部分岗位无法登陆，是否继续?', {
        btn: ['确定', '取消']
    }, function(){
        $.postJSON('/system/delrole', {
            Rid: postId
        }, function(res){
            if(res.ret == 200){   //成功
                layer.closeAll();
                getpage('/system/rolelist');
            }else {
                alert(res.msg);
            }
        }, 'post');
    });
});

