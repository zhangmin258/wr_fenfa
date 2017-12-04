
$(function(){

    //获取手机验证码
    var countdown = 60;
    function settime(obj){
        if (countdown == 0) {
            obj.removeAttr("disabled");
            obj.val("获取验证码");
            countdown = 60;
            return;
        }else{
            obj.attr("disabled", true);
            obj.val("重新发送(" + countdown + ")");
            countdown --;
        }
        setTimeout(function(){
            settime(obj);
        },1000);
    }

    //发送验证码
    $('#get-vcode').on('click', function(){
        var _this = $(this);
        var phone = $.trim($('#phone').val());
        //手机号 
        if(!phone){
            formError.show('手机号不得为空！');
            return;
        }
        if(phone && !/^1\d{10}$/.test(phone)){
            formError.show('手机号格式不正确！');
            return;
        }
        $.ajax({
            url: 'existaccount',
            type: 'post',
            data: {
                phone: phone
            },
            beforeSend: function(request) {
                request.setRequestHeader("X-Xsrftoken", $('meta[name=_xsrf]').attr('content'));
            },
            success: function(res){
                if(res.ret == 200){
                    formError.hide();
                    settime(_this);
                }else{
                    formError.show(res.err);
                }
            }
        });
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

    var page = {
        init: function(){
            this.bindEvent();
        },
        bindEvent: function(){
            var _this = this;
            //注册按钮点击事件
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
                Account: $.trim($('#phone').val()),
                Vcode: $.trim($('#v-code').val()),
                Password: $.trim($('#password').val()),
                confirmPassword: $.trim($('#confirm-password').val())
            }
            //表单验证结果
            var validateResult = this.formValidate(formDate);
            if(validateResult.status){   //如果验证通过
                //提交表单
                $.ajax({
                    url: '../account/register',
                    type: 'post',
                    data: formDate,
                    beforeSend: function(request) {  
                        request.setRequestHeader("X-Xsrftoken", $('meta[name=_xsrf]').attr('content'));  
                    },
                    success: function(res){
                        if(res.ret == 200){
                            layer.msg('注册成功！',{
                                shade: [0.3,'black'],
                                icon: 1,
                            });
                            setTimeout(function(){
                                window.location.href='login';
                            }, 1000);
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
            //手机号 
            if(!formDate.Account){
                result.msg = '手机号不得为空！';
                return result;
            }
            if(formDate.Account && !/^1\d{10}$/.test(formDate.Account)){
                result.msg = '手机号格式不正确！';
                return result;
            }

            //验证码
            if(!formDate.Vcode){
                result.msg = '验证码不得为空！';
                return result;
            }

            //密码
            if(!formDate.Password){
                result.msg = '密码不得为空！';
                return result;
            }
            if(formDate.Password && formDate.Password.length < 6){
                result.msg = '密码不得少于6位！';
                return result;
            }

            //确认密码
            if(!formDate.confirmPassword){
                result.msg = '确认密码不得为空！';
                return result;
            }
            if(formDate.confirmPassword && formDate.confirmPassword !== formDate.Password){
                result.msg = '2次输入的密码不一致！';
                return result;
            }

            //如果验证通过
            result.msg = '验证通过！';
            result.status = true;
            return result;
        }
    }

    page.init();

});
