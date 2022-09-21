function show_tab_content() {
    let current_pane = window.location.href.substring(window.location.href.lastIndexOf('/') + 1).substring(1);
    current_pane = current_pane === "" ? "home" : current_pane
    const query = window.location.search
    current_pane = (query.includes("code=") && query.includes("state=")) ? "home" : current_pane
    $(".tab-pane").each(function () {
        let dom_pane_id = $(this).attr('id');
        if ("tab-content-" + current_pane === dom_pane_id) {
            document.getElementById(dom_pane_id).style.display = "";
        } else {
            document.getElementById(dom_pane_id).style.display = "none";
        }
    })
}

$('#nav-tabs a').click(function (e) {
    // e.preventDefault();
    $(this).tab('show');
    show_tab_content();
    let current_pane = window.location.href.substring(window.location.href.lastIndexOf('/') + 1).substring(1);
    if (current_pane === "history") {
        // TODO: Find a way to just reload table width
        location.reload();
    }
});

// store the currently selected tab in the hash value
$("ul.nav-tabs > li > a").on("shown.bs.tab", function (e) {
    window.location.hash = $(e.target).attr("href").substring(1);
});

// on load of the page: switch to the currently selected tab
const hash = window.location.hash;
$('#nav-tabs a[href="' + hash + '"]').tab('show');
show_tab_content();
