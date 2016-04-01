// Data Pane

function toggleCommitData(e) {
    var parent = $(e.target).parent();
    var commitID = parent.parent().text().split(parent.text())[0].trim();
    var hidden = parent.find("span.glyphicon-minus").hasClass("hidden");

    parent.find("span").toggleClass("hidden");
    
    if ( hidden == true ) {
        $("tr." + commitID).removeClass("hidden");
    } else {
        $("tr." + commitID).addClass("hidden");
    }
}

function addDataPaneListeners() {
    $(".data .expand").off();
    $(".data .expand").click(toggleCommitData);
}


function displayPipelineOutput(result) {
    $(".output.data table").append(result);
    addDataPaneListeners();
}

// Poller to check if pipeline is done

function updatePipelineStatusUI(result) {
    if (pipelineCompleted !== null && pipelineCompleted) {
        return
    }

    pipelineCompleted = result["status"];

    if (pipelineCompleted) {
        loadPipelineOutput();

        $(".flash .message.info")
            .removeClass("hidden")
            .text("All pipelines have completed")
            .fadeOut(3000, function() { $(".flash .message.info").addClass("hidden") });

        window.clearInterval(pipelineStatusPoller);
    }

    $(".pane.code").removeClass("disabled");

    for(var outputRepo in result["states"]) {
        var statusBar = $(".status .item.hidden").clone().removeClass("hidden");
        var state = result["states"][outputRepo];

        statusBar.text( outputRepo + " is " + state );
        statusBar.appendTo(".status");
        fade(statusBar);
    }

}

function fade(elem) {
    var fade_helper = function() {
        elem.fadeOut(2000);
    }
    window.setTimeout(fade_helper, 2000);
}

function checkPipelineStatus() {

    $(".pane.code").addClass("disabled");

    $.ajax({
            url: "/pipeline/status",
            success: updatePipelineStatusUI
    });
}

function loadPipelineOutput() {
    $.ajax({
            url: "/pipeline/output",
                success: displayPipelineOutput
        });
}

function enableRunButton() {
    $("button.run").on("click", function() {
            $("form").submit();
        });
}

// Initialization

var pipelineStatusPoller = window.setInterval(checkPipelineStatus, 250);
var pipelineCompleted = null;

function initializeCodeMirror() {
    var myTextArea = $("textarea[name='code']")[0];
    var myCodeMirror = CodeMirror.fromTextArea(myTextArea);
    myCodeMirror.setOption("mode",{name:"javscript", json: true});
    myCodeMirror.setOption("lineNumbers", true);
    
}

$(document).ready(
                  function () {
                      initializeCodeMirror();
                      enableRunButton();
                      addDataPaneListeners();
                      $(".steps").on('afterChange', function(event, slick, currentSlide){
                              $("input[name='current_step']").attr("value",currentSlide);
                          });
                      var initialSlide = $(".steps").attr("current-slide");
                      $(".steps").slick({
                                  appendArrows: $(".arrows"),
                                  slidesToShow:1, 
                                  slidesToScroll:1, 
                                  infinite: false,
                                  initialSlide: parseInt(initialSlide),
                                  prevArrow: '<button type="button" class="btn btn-default slick-prev">Previous</button>',
                                  nextArrow: '<button type="button" class="btn btn-default slick-next">Next</button>',
                          });                      

                      $(".next").on("click", function() { $(".slick-next").click() });
                      $(".prev").on("click", function() { $(".slick-prev").click() });

                  }
                  );