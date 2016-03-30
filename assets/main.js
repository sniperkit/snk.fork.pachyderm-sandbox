// Poller to check if pipeline is done

function updatePipelineStatusUI(result) {
    if (pipelineCompleted !== null && pipelineCompleted) {
        return
    }

    console.log(result);
    console.log(result["status"]);
    pipelineCompleted = result["status"];

    if (pipelineCompleted) {
        $(".flash .message").text("All pipelines have completed");
        $(".flash").removeClass("hidden");
    }

    $(".pane.code").removeClass("disabled");

    for(var outputRepo in result["states"]) {
        var statusBar = $(".status .item.hidden").clone().removeClass("hidden");
        
        for(var commitID in result["states"][outputRepo]) {
            var state = result["states"][outputRepo][commitID];
            statusBar.text( outputRepo + " completed commit (" + commitID + ")" + state );
            statusBar.appendTo(".status");
            statusBar.fadeOut(2000);
        }
    }

    window.clearInterval(pipelineStatusPoller);
}

function checkPipelineStatus() {

    $(".pane.code").addClass("disabled");

    $.ajax({
            url: "/check_pipeline_status",
            success: updatePipelineStatusUI
    });
}

var pipelineStatusPoller = window.setInterval(checkPipelineStatus, 250);

var pipelineCompleted = null;