
$(function(){
    //表单的错误提示信息
    var formError = {
        show: function(errMsg){
            $('.error-item').show().find('.error-msg').text(errMsg);
        },
        hide: function(){
            $('.error-item').hide().find('.error-msg').text('');
        }
    }

    var page = {
        init: function(){
            this.bindEvent();
        },
        bindEvent: function(){
            var _this = this;

            //登录按钮点击事件
            $('#submit').click(function(){
                _this.submit();
            });
            //当按了回车键
            $('.user-content').on('keyup', function(e){
                if(e.keyCode === 13){
                     _this.submit();
                }
            });
        },
        //提交表单
        submit: function(){
            var formDate = {
                account: $.trim($('#username').val()),
                password: $.trim($('#password').val())
            }
            //表单验证结果
            var validateResult = this.formValidate(formDate);
            if(validateResult.status){   //如果验证通过
                //提交表单
                $.ajax({
                    url: '/login',
                    type: 'post',
                    data: formDate,
                    beforeSend: function(request) {  
                        request.setRequestHeader("X-Xsrftoken", $('meta[name=_xsrf]').attr('content'));  
                    },
                    success: function(res){
                        if(res.ret == 200){
                            layer.msg(res.msg, {
                                shade: [0.3,'black'],
                                icon: 1
                            });
                            if(res.jumpBand === 1 && res.userType === 0){
                                window.location.href = '/usersbankcard/bindcardpage';   //绑卡
                            }else{
                                window.location.href = '/data/datapage';    //数据统计
                            }
                        }else{
                            formError.show(res.err);
                        }
                    }
                });
            }else{    //如果验证不通过
                formError.show(validateResult.msg);
            }
        },
        //表单字段验证
        formValidate: function(formDate){
            var result = {
                status: false,
                msg: ''
            }
            //如果用户名为空  
            if(!formDate.account){
                result.msg = '用户名不得为空！';
                return result;
            }
            //如果密码为空
            if(!formDate.password){
                result.msg = '密码不得为空！';
                return result;
            }
            //如果验证通过
            if(formDate.account && formDate.password){
                result.msg = '验证通过！';
                result.status = true;
                return result;
            }
        }
    }

    page.init();
});