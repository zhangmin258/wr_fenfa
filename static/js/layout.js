
$(function(){
	jQuery('.main-content').css({'min-height':$(window).height()});
	// $.ajax({
	// 	url: '/',
	// 	type: 'post',
 //        contentType:'application/json;charset=utf-8',
 //        dataType: 'text',
	// 	success: function(result){
 //            if (result && result.ret == 200 && !!result.data){
 //                var menu = '';
 //                var count = 1;
 //                $.each(result.data,function(i,m){
 //                    if (m.ChildMenu && m.ChildMenu.length > 0) {
 //                        menu+='<li class="menu-list menuList'+ count +'"><a href="#" ><span>'+ m.TabName+'</span></a><ul class="sub-menu-list">';
 //                        count ++ ;
 //                        var subMenu='';
 //                        $.each(m.ChildMenu,function(j,n){
 //                            subMenu+='<li><a class="norefresh" href="'+n.LinkUrl+'"> '+n.TabName+'</a></li>';
 //                        });
 //                        menu+=subMenu+'</ul></li>';
 //                    }else{
 //                        menu+= '<li><a class="norefresh" href="'+m.LinkUrl+'" ><span>'+ m.TabName+'</span></a></li>';
 //                    }
 //                });
 //                $('.js-left-nav').append(menu);
 //                leftSelect();
 //            }
	// 	},
	// 	error: function () {
	// 		alert('错误！');
 //        }
	// });
    
    $.postJSON('/', {}, function(result){
        if (result && result.ret == 200 && !!result.data){
            var menu = '';
            var count = 1;
            $.each(result.data,function(i,m){
                if (m.ChildMenu && m.ChildMenu.length > 0) {
                    menu+='<li class="menu-list menuList'+ count +'"><a href="#" ><span>'+ m.TabName+'</span></a><ul class="sub-menu-list">';
                    count ++ ;
                    var subMenu='';
                    $.each(m.ChildMenu,function(j,n){
                        subMenu+='<li><a class="norefresh" href="'+n.LinkUrl+'"> '+n.TabName+'</a></li>';
                    });
                    menu+=subMenu+'</ul></li>';
                }else{
                    menu+= '<li><a class="norefresh" href="'+m.LinkUrl+'" ><span>'+ m.TabName+'</span></a></li>';
                }
            });
            $('.js-left-nav').append(menu);
            leftSelect();
        }else{
            alert('错误！');
        }
    }, 'post');




	//左边菜单加选中状态
	function leftSelect(){
		var pathname = location.pathname;
        if (pathname!="/"){
        	$('.js-left-nav .norefresh').filter(function(){   
	        	return $(this).attr('href') == pathname;
	        }).closest("li").addClass('active').parents('.menu-list').addClass('nav-active');
        }
	}

    $("body").delegate(".menu-list > a","click",function(){
        var parent = jQuery(this).parent();
        var sub = parent.find('> ul');
        //if(!jQuery('body').hasClass('left-side-collapsed')) {
         	if(sub.is(':visible')) {
         		parent.removeClass('nav-active');
	            // sub.slideUp(200, function(){
	            //    parent.removeClass('nav-active');
	            //    jQuery('.main-content').css({height: ''});
	            //    mainContentHeightAdjust();
	            // });
         	} else {
	            visibleSubMenuClose();
	            parent.addClass('nav-active');
	            sub.slideDown(200, function(){
	                //mainContentHeightAdjust();
            	});
        	}
        //}
      	return false;
    });

    var firstEnter = false;
	$("body").delegate(".js-left-nav .norefresh","click",function(){
		$('.js-left-nav .active').removeClass('active');
		$(this).closest("li").addClass('active');
		if(!jQuery('body').hasClass('left-side-collapsed') && $('.js-left-nav .nav-active').length && $(this).closest(".nav-active").length==0) {
			visibleSubMenuClose();
		}
		var url = $(this).attr('href');
		getpage(url, firstEnter);
		firstEnter = true;
		// $('.wrapper').zload(url);
      	return false;
    });


	window.onpopstate = function (e) {
		if (e.state){
			$('.js-left-nav .active').removeClass('active');
			 visibleSubMenuClose();
			leftSelect();
			$('.wrapper').empty().append(e.state.html);
			execjs(e.state.html);
		}
	}

    function visibleSubMenuClose() {
      	jQuery('.menu-list').each(function() {
         	var t = jQuery(this);
         	if(t.hasClass('nav-active')) {
	            t.find('> ul').slideUp(200, function(){
	               t.removeClass('nav-active');
	            });
         	}
      	});
    }
});

function getpage(url,isFirst){
	$.zget(url,"",function(result){
		history.pushState({html:result}, "what", url);
        $('.wrapper').empty().append(result);
        if(!isFirst){
			execjs(result);
        }
	});
	return false;
}


$('html').on('click','.skip1',function () {
	
	window.sessionStorage.URl = window.location.href;

    var phone = $.trim($("#phone").val());   //手机号码
    window.sessionStorage.phone = phone;

    var href = $(this).attr('href');
	getpage(href);
	return false;
});



$('html').on('click','.skip2',function () {

	window.sessionStorage.URl = window.location.href;

    // var productName = $.trim($('#product-name').val());   //产品名称
    // window.sessionStorage.productName = productName;

	return false;
});



$('html').on('mousedown','.norefresh',function () {
    window.sessionStorage.clear();

});


// $('html').on('mouseup','.pull-left button[type="submit"]',function () {
//     window.sessionStorage.clear();
// })

// function zaq(name) {
//     var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
//     var r = window.location.search.substr(1).match(reg);
//     if(r!=null)return  unescape(r[2]); return null;
// }


function execjs(html){
	// 第一步：匹配加载的页面中是否含有js
	var regDetectJs = /<script(.|\n)*?>(.|\n|\r\n)*?<\/script>/ig;
	var jsContained = html.match(regDetectJs);
	// 第二步：如果包含js，则一段一段的取出js再加载执行
	if(jsContained) {
		// 分段取出js正则
		var regGetJS = /<script(.|\n)*?>((.|\n|\r\n)*)?<\/script>/im;

		// 按顺序分段执行js
		var jsNums = jsContained.length;
		for (var i=0; i<jsNums; i++) {
			var jsSection = jsContained[i].match(regGetJS);
			if(jsSection[2]) {
				if(window.execScript) {
					// 给IE的特殊待遇
					window.execScript(jsSection[2]);
				} else {
					// 给其他大部分浏览器用的
					window.eval(jsSection[2]);
				}
			}
		}
	}
}
