$(function () {
    $('#logout').on('click', function () {
        if (confirm("退出登录？")) {
            window.location.href = 'logout';
        }
    });
    $('input').keypress(function (e) {
        if (e.which == 13) {
            sendMsg();
            $('#content').val('');
        }
    });
  	$('input').on('focus', function(){
      $('html, body').animate({scrollTop: $(document).height()}, 'fast');
    });
    $('input').on('blur', function(){
    });
    $('#myForm button').on('click', function (e) {
        e.preventDefault();
        sendMsg();
    });
})
function edit() {
    window.location.href = 'edit';
}
function sendMsg() {
    $.ajax({
        url: 'post',
        type: 'post',
        dataType: 'json',
        data: $('#myForm').serialize(),
        success: function (data) {

            $('#content').val('');
        }
    });
}
var lastReceived = 0;
var isWait = false;

var fetch = function () {
    if (isWait) return;
    isWait = true;
    $.getJSON("/fetch?lastReceived=" + lastReceived, function (data) {
        if (data == null) return;
        console.log(data);
        $.each(data, function (i, event) {
            // var li = document.createElement('li');
            switch (event.Type) {
                case 0: // JOIN
                    joinAdd(event.Username);
                    break;
                case 1: // LEAVE
                    leaveAdd(event.Username);
                    break;
                case 2: // MESSAGE
                    joinMsg(event);
                    break;
            }
            $(document).scrollTop($(document).height());

            // $('#chatbox li').first().before(li);

            lastReceived = event.Addtime;
        });
        isWait = false;
    });
}

// Call fetch every 3 seconds
setInterval(fetch, 3000);

fetch();

function joinAdd(username) {
    var myDiv = $('<div></div>');
    myDiv.addClass('join');
    myDiv.html(username + ' 加入聊天');
    $('main').append(myDiv);
}

function leaveAdd(username) {
    var myDiv = $('<div></div>');
    myDiv.addClass('join');
    myDiv.html(username + ' 离开聊天');
    $('main').append(myDiv);
}

function joinMsg(data) {
    var msg = $('<div></div>');
    msg.addClass('msg');
    var a = $('<a href="javascript:;" /></a>');
        
    msg.append(a);

    if (data.ID == $('#user_id').val()) {
        msg.addClass('msg_self');
        a.attr('href', "javascript:edit();");
    }

    var img = $('<img src="" />');
    img.addClass('head');
    if (data.Headimgurl == '') {
        img.attr('src', "/static/home/images/dog.jpg");
    } else {
        img.attr('src', data.Headimgurl);
    }
    a.append(img);
    

    var msg_box = $('<div></div>');
    msg_box.addClass('msg_box');
    msg.append(msg_box);

    var msg_top = $('<div></div>');
    msg_top.addClass('msg_top');
    msg_top.html(data.Username);
    var msg_span = $('<span></span>');
    msg_span.html(timestampToTime(data.Addtime));
    if (data.ID == $('#user_id').val()) {
        msg_top.prepend(msg_span);
    } else {
        msg_top.append(msg_span);
    }


    msg_box.append(msg_top);


    var msg_bottom = $('<div></div>');
    msg_bottom.addClass('msg_bottom');
    msg_box.append(msg_bottom);

    var msg_i = $('<i></i>');
    msg_bottom.append(msg_i);

    var msg_content = $('<div></div>');
    msg_content.addClass('msg_content');
    msg_content.html(data.Content);
    msg_bottom.append(msg_content);

    $('main').append(msg);
}
/*
* 时间戳转日期
* @param timestamp
* @returns {*}
*/
function timestampToTime(timestamp) {
    var nowDate = new Date();
    var date = new Date(timestamp * 1000);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
    Y = date.getFullYear() == nowDate.getFullYear() ? '' : date.getFullYear() + '-';

    M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + '-';
    if (Y == '') {
        if (date.getMonth() == nowDate.getMonth()) {
            M = '';
        }
    }
    D = date.getDate() + ' ';
    if (Y == '' && M == '') {
        if (date.getDate() == nowDate.getDate()) {
            D = '';
        }
    }
    h = date.getHours() + ':';
    m = (date.getMinutes() < 10 ? '0' + (date.getMinutes()) : date.getMinutes()) + ':';
    s = (date.getSeconds() < 10 ? '0' + (date.getSeconds()) : date.getSeconds());

    return Y + M + D + h + m + s;
}