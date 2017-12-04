
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
                UserName: $.trim($('#user').val()),  //姓名
                IdCard: $.trim($('#ID').val()),     //身份证号
                BankCardNumber: $.trim($('#bankCard-number').val()),  //银行卡号
                BankMobile: $.trim($('#bank-phone').val())    //银行预留手机号
            }
            //表单验证结果
            var validateResult = this.formValidate(formDate);
            if(validateResult.status){   //如果验证通过
                //提交表单
                $.postJSON('../usersbankcard/bindcard', formDate, function(res){
                    if(res.ret == 200){
                        layer.msg('银行卡绑定成功！',{
                            shade: [0.3,'black'],
                            icon: 1
                        });
                        window.location.href = '/data/datapage';    //数据统计
                    }else{
                        formError.show(res.err);
                    }
                }, 'post');
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
            //姓名
            if(!formDate.UserName){
                result.msg = '姓名不得为空！';
                return result;
            }
            //身份证号
            if(!formDate.IdCard){
                result.msg = '身份证号不得为空！';
                return result;
            }
            //15位数的身份证号码验证
            if(formDate.IdCard && formDate.IdCard.length <= 15 && !/^[1-9]\d{7}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}$/.test(formDate.IdCard)){
                result.msg = '身份证号格式不正确！';
                return result;
            }
            //18位数的身份证号码验证
            if(formDate.IdCard && formDate.IdCard.length <= 18 && !/^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}([0-9]|X)$/.test(formDate.IdCard)){
                result.msg = '身份证号格式不正确！';
                return result;
            }
            //银行卡号
            if(!formDate.BankCardNumber){
                result.msg = '银行卡号不得为空！';
                return result;
            }
           /* if(formDate.BankCardNumber && !/^(\d{16}|\d{19})$/.test(formDate.BankCardNumber)){
                result.msg = '银行卡号格式不正确！';
                return result;
            }*/
            //银行预留手机号
            if(!formDate.BankMobile){
                result.msg = '手机号不得为空！';
                return result;
            }
            if(formDate.BankMobile && !/^1\d{10}$/.test(formDate.BankMobile)){
                result.msg = '手机号格式不正确！';
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
