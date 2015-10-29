$(document).ready(function() {

    function displaySitemap(nodes) {

        var $sitemap = $("#sitemap");
        $sitemap.empty();

        for (var i = 0; i < nodes.length; i++) {
            var node = nodes[i];

            var $article = $("<article></article>");

            var $header = $("<header></header>");
            if (node.offsite) {
                $header.addClass("offsite");
            }
            var $url = $("<a></a>");
            $url.attr("target", "_blank");
            $url.attr("href", node.url);
            $url.text(node.url);
            $header.append($url);

            var $list = $("<ul></ul>");
            $list.addClass("list-unstyled");
            for (var j = 0; j < node.links.length; j++) {
                var $elem = $("<li></li>");
                var $link = $("<a></a>");
                $link.attr("target", "_blank");
                $link.attr("href", node.links[j]);
                $link.text(node.links[j]);

                $elem.append($link);
                $list.append($elem);
            }
            if (node.offsite || node.links.length == 0) {
                var $elem = $("<li></li>");
                $elem.text(node.offsite ? "Not within domain so did not crawl." : "No links found.");
                $list.append($elem);
            }
            $article.append($header).append($list);

            $sitemap.append($article);
        }
    }
    
    $(".loading").hide();

    $("#crawl").submit(function(e) {
        e.preventDefault();

        var $this = $(this);
        var $loading = $this.find(".loading");
        var $original = $this.find(".original");
        var url = $this.find("#url").val();

        $loading.show();
        $original.hide();

        $.get("/crawl", {url: url})
            .done(function(data) {
                displaySitemap(data);
            })
            .fail(function(data) {
                alert("Error (" + data.status + "): " + data.responseText);
            })
            .always(function() {
                $loading.hide();
                $original.show();
            });
    });

});
