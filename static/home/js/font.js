// 100px = 1rem
!function (n) {
    // window.document 指向当前窗口内的文档节点
    // alert( window.document.title);    // 弹出: title
    var e = n.document,
        t = e.documentElement,
        i = 720,
        d = i / 100,
        // 判断屏幕是否有调整变化
        // 如果window没有orientationchange就用resize
        o = "orientationchange" in n ? "orientationchange" : "resize",
        a = function () {
            // 如果没有clientWidth则取320
            // 如果大于720则取720
            // 以屏幕320像素为基准
            var n = t.clientWidth || 320; n > 720 && (n = 720);
            t.style.fontSize = n / d + "px"
        };
    e.addEventListener && (n.addEventListener(o, a, !1), e.addEventListener("DOMContentLoaded", a, !1))
}(window);