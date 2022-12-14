function getQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return unescape(r[2]);
    return null;
}


$('#checkboxall').click(function () {
    if ($('#checkboxall').prop('checked') == true) {
        $('.itemcheck').prop('checked', true)
    } else {
        $('.itemcheck').prop('checked', false)
    }
});

$('#checkboxrev').click(function () {
    $.each($('.itemcheck'), function (e) {
        if ($(this).prop('checked') == true) {
            $(this).prop('checked', false)
        } else {
            $(this).prop('checked', true)
        }
    })
});

$('#field-type').change(function () {
    type = $('#field-type').val();
    risk = $('#field-risk').val();
    search = $('#field-search').val();
    $('#field-plugin').children().remove();
    $.post('/plugin', {
        type: type,
        risk: risk,
        search: search
    }, function (e) {
        $.each(e, function (i, n) {
            $('#field-plugin').append("<option title='" + n['info'] + "'>" + n['name'] + "</option>");
        });
        $('#field-plugin').multiSelect('refresh');
    }, "json")
});

$('#field-risk').change(function () {
    type = $('#field-type').val();
    risk = $('#field-risk').val();
    search = $('#field-search').val();
    $('#field-plugin').children().remove();
    $.post('/plugin', {
        type: type,
        risk: risk,
        search: search
    }, function (e) {
        $.each(e, function (i, n) {
            $('#field-plugin').append("<option title='" + n['info'] + "'>" + n['name'] + "</option>");
        });
        $('#field-plugin').multiSelect('refresh');
    }, "json")
});

$('#field-search').bind('keyup', function (event) {
    if (event.keyCode == "13") {
        type = $('#field-type').val();
        risk = $('#field-risk').val();
        search = $('#field-search').val();
        $('#field-plugin').children().remove();
        $.post('/plugin', {
            type: type,
            risk: risk,
            search: search
        }, function (e) {
            $.each(e, function (i, n) {
                $('#field-plugin').append("<option title='" + n['info'] + "'>" + n['name'] + "</option>");
            });
            $('#field-plugin').multiSelect('refresh');
        }, "json")
    }
});

$('#savetask').click(function () {
    title = $('#field-name').val();
    if (!title) {
        swal("????????????????????????", "", "error");
    }
    condition = getQueryString('q');
    plugin = $('#field-plugin').val().join(",");
    plan = $('#field-plan').val();
    isupdate = $('#field-isupdate').val() == "" ? "0" : $('#field-isupdate').val();
    resultcheck = $('#resultcheck').prop('checked')
    var ids = [];
    if (!resultcheck) {
        $.each($('.itemcheck:checked'), function (i, n) {
            ids.push($(n).attr('infoid'))
        });
        ids = ids.join(',');
    }
    $.post('/task', {
        title: title,
        condition: condition,
        plugin: plugin,
        ids: ids,
        plan: plan,
        isupdate: isupdate,
        resultcheck: resultcheck
    }, function (e) {
        if (e == 'success') {
            swal("????????????", '', "success");
            $('.confirm').click(function () {
                location.href = "/task";
            })
        } else {
            swal("????????????", "??????????????????????????????????????????????????????!", "error")
        }
    })
});

$('#field-plan').change(function () {
    if ($(this).val() == 0) {
        $('#field-isupdate').attr('disabled', 'disabled')
    } else {
        $('#field-isupdate').removeAttr('disabled')
    }
});

$('#addnewitems').click(function () {
    ips = $('#field-newitems').val().split(',');
    for (i in ips) {
        $('#content').append("<div class='col-lg-4'>\
                                <div class='card-box project-box'>\
                                    <a href='javascript:;' style='position: absolute;right:10px'>\
                                        <input type='checkbox' class='itemcheck' infoid='" + ips[i] + "'>\
                                    </a>\
                                    <h4 class='m-t-0 m-b-5' style='height: 25px'>\
                                    <a href='http://" + ips[i] + "' class='text-name' target='_blank'>" + ips[i] + "</a>\
                                    </h4>\
                                    <div class='clearfix'></div>\
                                    </div>\
                                    </div>");
    }
    $('#tips').html("");
    $('#closenewitem').click();
});

$('.tag').click(function () {
    if ($(this).hasClass('zmdi-plus')) {
        $(this).removeClass('zmdi-plus').addClass('zmdi-minus')
    } else {
        $(this).removeClass('zmdi-minus').addClass('zmdi-plus')
    }
});

$(document).ready(function () {
    if ($('#content').html().trim() == '') {
        if (getQueryString('q') == null) {//???????????????
            $('.wrapper .container').append("<div style='font-size:22px;' id='tips'>????????????????????????</div>")
        } else {//???????????????
            $('.wrapper .container').append("<div style='font-size:22px;' id='tips'>??????????????????????????????</div>")
        }
    }
});

$('#select-all').click(function(){
    if($('#field-plugin').val()==null){
        $('#field-plugin').multiSelect('select_all');
    }else{
        $('#field-plugin').multiSelect('deselect_all');
    }
});


