
$(document).ready(function($) {

    $("#backtotop").hide();
    var isVisible = false;
    $(window).scroll(function(){
        var pageBottom = $(window).scrollTop() + $(window).height();
        var shouldBeVisible = $(window).scrollTop()>200;
        if (shouldBeVisible && !isVisible || pageBottom === $(document).height()) {
            isVisible = true;
            $("#backtotop").show();
            $("#scroll-down").hide();
        } else if (isVisible && !shouldBeVisible && pageBottom !== $(document).height()) {
            isVisible = false;
            $("#backtotop").hide();
        }
    });

    var container = $("#viewer");

    $(".pictures").last().clone().appendTo(container);

    var clickEvent = (function() {
        if ('ontouchstart' in document.documentElement === true)
            return 'touchstart';
        else
            return 'click';
    })();

    $("img").on(clickEvent, function(e) {
        $("#viewer").empty();
        $(this).clone().appendTo($("#viewer").css("visibility", "visible"));
    });

    $(document).mouseup(function(e) {

        if (!$("#viewer").is(e.target) && $("#viewer").has(e.target).length === 0 && $("#viewer").css("visibility") == "visible")
        {
            $("#viewer").empty();
            $("#viewer").css("visibility", "hidden");
        }
    });

    var w = 10;

    $("input[name='name'], input[name='email'], input[name='tel']").change(function () {
        var prog = $(".progress_bar");

        if ($(this).is(":valid")) {
            w += 10;
        } else if (w >= 10) {
            w -= 10;
        }
        $("input[name='progress_width']").val(w);
        prog.animate({width: w + "%"}, "slow");
    });

    $("#service-form").submit(function(event) {
        $(".wait").addClass("spinner");
    });

    if ($(".thanks").length) {
        $(this).delay(5000).queue(function () {
            window.location.replace("/");
        });
    }

    var flips = [false, false, false];

    $(".flip-box").on(clickEvent, function(e) {
        var i = $(".flip-box").index(this);
        if (!flips[i]) {
            $(this).find(".flip-box-inner").addClass("jqueryClickFlip");
            $(this).find(".flip-box-front").addClass("jqueryClickHide");
            flips[i] = true;
        } else {
            $(this).find(".flip-box-inner").removeClass("jqueryClickFlip");
            $(this).find(".flip-box-front").removeClass("jqueryClickHide");
            flips[i] = false;
        }
    });

    $("input[type='submit']").click(function() {


        $("#service-form-parent").effect("shake");
    });

    //thumbnail slider ------------------------------>
    $(".show-slider-svg").on(clickEvent, function () {
        $("#thumbnail-slider").removeClass("hide-slider");
        $("#thumbnail-slider").addClass("show-slider");
        $(".hide-slider-svg").removeClass("svg-selected");
        $(this).addClass("svg-selected");
        $("#pic-links a").css("height", "50px");
        $("#thumbnail-slider label").show();
        $("#pic-links a").css("border", "2px solid transparent");
        var hash = window.location.hash;
        $("#pic-links a[href='" + hash + "']").css("border", "2px solid rgb(30,124,253)");
    });
    $(".hide-slider-svg").on(clickEvent, function () {
        $("#thumbnail-slider").removeClass("show-slider");
        $("#thumbnail-slider").addClass("hide-slider");
        $(".show-slider-svg").removeClass("svg-selected");
        $(this).addClass("svg-selected");
        $("#pic-links a").css("height", "0px");
        $("#thumbnail-slider label").hide();
        $("#pic-links a").css("border", "none");
    });

    function rightEnd() {
        $("label[for='slide-right']").addClass("disable-right-arrow");
        $("label[for='slide-left']").removeClass("disable-left-arrow");
        $("label[for='slide-left']").addClass("left-side-hover");
        $("label[for='slide-right']").removeClass("right-side-hover");
    }

    function leftEnd() {
        $("label[for='slide-left']").addClass("disable-left-arrow");
        $("label[for='slide-right']").removeClass("disable-right-arrow");
        $("label[for='slide-right']").addClass("right-side-hover");
        $("label[for='slide-left']").removeClass("left-side-hover");
    }

    var thumb_count = $("#pic-links a").length;
    var thumb_width = $("#pic-links a").height();
    var window_width = $("#pic-links").width();
    var slider_width = thumb_count * thumb_width;
    var thumbs_in_view = window_width / thumb_width;

    var view_num = 1;
    var view_count = Math.ceil(slider_width / (thumbs_in_view * thumb_width));

    var current_slide;

    var view_start = 0;
    var view_end = view_num * window_width;

    function slide() {
        current_slide = parseInt($("#slide-num").val());

        var spot = current_slide * thumb_width;

        view_end = view_num * window_width;

        for (var i = 1; i <= view_count; i++) {
            if (spot <= i * window_width) {
                //alert("view num: " + i);
                view_num = i;
                break;
            }
        }

        //alert(spot + ", " + view_end);

        if (spot > view_end) {

            view_start = view_end;

            view_end = view_num * window_width;

            //alert("view num: " + view_num + ", " + "view start: " + view_start);

            $("#pic-links").animate({ scrollLeft: view_start}, 650);
            if (view_end >= slider_width) {
                rightEnd();
            }
        }
        else if (spot < view_start) {

            //alert(spot + " < " + view_start);

            view_end = view_num * window_width;

            view_start = view_end - window_width;

            //alert(view_start + " -> " + view_end);

            $("#pic-links").animate({ scrollLeft: view_start}, 650);
            if (view_start === 0) {
                leftEnd();
            }
        }
    }

    slide();

    $("#slide-num").change("input", function() {
        slide();
    });

    $("label[for='slide-right']").on(clickEvent, function() {
        //alert(view_end + ", " + slider_width);
        if (view_end < slider_width) {
            view_start += window_width;
            view_end += window_width;
            $("#pic-links").animate({ scrollLeft: view_start}, 650);
        }
    });
    $("label[for='slide-left']").on(clickEvent, function() {
        if (view_start !== 0) {
            view_start -= window_width;
            view_end -= window_width;
            $("#pic-links").animate({ scrollLeft: view_start}, 650);
        }
    });

    $(document).keydown(function(e) {
        switch(e.which) {
            case 37: // left
                window.location = $("#pic-left").attr("href");
                break;

            // case 38: // up
            //     break;

            case 39: // right
                window.location = $("#pic-right").attr("href");
                break;

            // case 40: // down
            //     break;

            default: return; // exit this handler for other keys
        }
        e.preventDefault(); // prevent the default action (scroll / move caret)
    });
    //end thumbnail slider ------------------------------>

    $(document).on('click', '.focus-select', function(){
        // $(this).val("P");
    });

    $(".form-inputs").focusout(function() {
        $(this).addClass("error-style");
    });

    /**
     * Copyright 2012, Digital Fusion
     * Licensed under the MIT license.
     * http://teamdf.com/jquery-plugins/license/
     *
     * @author Sam Sehnert
     * @desc A small plugin that checks whether elements are within
     *     the user visible viewport of a web browser.
     *     only accounts for vertical position, not horizontal.
     */

    $.fn.visible = function(partial) {

        var $t            = $(this),
            $w            = $(window),
            viewTop       = $w.scrollTop(),
            viewBottom    = viewTop + $w.height(),
            _top          = $t.offset().top,
            _bottom       = _top + $t.height(),
            compareTop    = partial === true ? _bottom : _top,
            compareBottom = partial === true ? _top : _bottom;

        return ((compareBottom <= viewBottom) && (compareTop >= viewTop));

    };

    $(window).scroll(function(event) {
        var el_1 = $("#mission");
        var el_2 = $("#headline");
        if (el_1.length && el_1.visible(true)) {
            el_1.addClass("mission-appear");
            el_2.addClass("headline-appear");
        }
        var el_3 = $(".services-content ol li");
            el_3.each(function () {
                if ($(this).visible(true)) {
                    $(this).addClass("services-appeared");
                }
            });
        var el_4 = $(".finance-content ul li");
        if (el_4.length && el_4.eq(0).visible(true)) {
            el_4.eq(0).addClass("finance-content-appeared-first-li");
        }
        if (el_4.length && el_4.eq(1).visible(true)) {
            el_4.eq(1).addClass("finance-content-appeared-second-li");
            $("#financiers a").addClass("financiers-appeared");
        }
        var el_5 = $(".card");
        el_5.each(function () {
            if ($(this).visible(true)) {
                $(this).addClass("card-appeared");
            }
        });
    });
});

var createdStyleTag = document.createElement("style");

var content_width = window.getComputedStyle(
    document.querySelector('.our-mother-earth'), '::after'
).getPropertyValue('width');

var keyFrames = ".our-mother-earth::after {\n" +
    "   animation: 8s linear 8s forwards text-ray-width;\n" +
    "   }\n" +
    "   @keyframes text-ray-width {\n" +
    "        0% {\n" +
    "            opacity: 0;\n" +
    "            width: 0;\n" +
    "        }\n" +
    "        100% {\n" +
    "            opacity: 1;\n" +
    "            width: A_DYNAMIC_VALUE;\n" +
    "        }\n" +
    "    };";

createdStyleTag.innerText = keyFrames.replace(/A_DYNAMIC_VALUE/g, content_width);
document.body.appendChild(createdStyleTag);








