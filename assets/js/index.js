$(document).ready(function() {
    
    $(".loading").hide();

    $("#crawl").submit(function(e) {
        e.preventDefault();

        var $this = $(this);
        var $loading = $this.find(".loading");
        var $original = $this.find(".original");

        $loading.show();
        $original.hide();
        setTimeout(function() {
            $loading.hide();
            $original.show();
        }, 2000);
    });

});
