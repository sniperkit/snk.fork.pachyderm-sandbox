// Poller to check if pipeline is done

function updatePipelineStatusUI(result) {
    console.log(result);
    console.log(result["status"]);

    $(".pane.code").removeClass("disabled");

    for(var outputRepo in result["states"]) {
        var statusBar = $(".status .item.hidden").clone().removeClass("hidden");
        
        for(var commitID in result["states"][outputRepo]) {
            var state = result["states"][outputRepo][commitID];
            statusBar.text( outputRepo + " completed commit (" + commitID + ")" + state );
            statusBar.appendTo(".status");
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