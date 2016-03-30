// Data Pane

function toggleCommitData(e) {
    var parent = $(e.target).parent();
    var commitID = parent.text().split($(e.target).text())[0].trim();

    if ( $(e.target).text() == "+") {
        $("tr." + commitID).removeClass("hidden");
        $(e.target).text("-");
    } else {
        $("tr." + commitID).addClass("hidden");
        $(e.target).text("+");
    }
}

function addDataPaneListeners() {
    $(".data .expand").click(toggleCommitData);
}


function displayPipelineOutput(result) {

}

// Poller to check if pipeline is done

function updatePipelineStatusUI(result) {
    if (pipelineCompleted !== null && pipelineCompleted) {
        return
    }

    console.log(result);
    console.log(result["status"]);
    pipelineCompleted = result["status"];

    if (pipelineCompleted) {
        $(".flash .message.info")
            .removeClass("hidden")
            .text("All pipelines have completed")
            .fadeOut(3000, function() { $(".flash .message.info").addClass("hidden") });
    }

    $(".pane.code").removeClass("disabled");

    for(var outputRepo in result["states"]) {
        var statusBar = $(".status .item.hidden").clone().removeClass("hidden");
        
        for(var commitID in result["states"][outputRepo]) {
            var state = result["states"][outputRepo][commitID];
            statusBar.text( outputRepo + " completed commit (" + commitID + ")" + state );
            statusBar.appendTo(".status");
            fade(statusBar);
        }
    }

    window.clearInterval(pipelineStatusPoller);
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

// Initialization

var pipelineStatusPoller = window.setInterval(checkPipelineStatus, 250);
var pipelineCompleted = null;

$(document).ready(
                  function () {
                      addDataPaneListeners();
                  }
                  );