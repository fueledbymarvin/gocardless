$(document).ready(function() {

    function drawGraph(graph) {

        var width = 960,
            height = 500;

        var force = d3.layout.force()
                .charge(-120)
                .linkDistance(30)
                .size([width, height]);

        d3.select("svg").remove();
        var svg = d3.select("section#graph").append("svg")
                .attr("width", width)
                .attr("height", height);

        force.nodes(graph.nodes).links(graph.links).start();

        var links = svg.selectAll("line")
                .data(graph.links)
                .enter()
                .append("line")
                .attr("class", "link");

        var nodes = svg.selectAll("g")
                .data(graph.nodes)
                .enter()
                .append("g");

        nodes.append("circle")
            .attr("class", "node")
            .attr("r", 5)
            .call(force.drag);

        nodes.append("text")
            .text(function(d) { console.log(d.url); return d.url; });

        force.on("tick", function() {
            links.attr("x1", function(d) { return d.source.x; })
                .attr("y1", function(d) { return d.source.y; })
                .attr("x2", function(d) { return d.target.x; })
                .attr("y2", function(d) { return d.target.y; });

            nodes.attr("transform", function (d) {
                return "translate(" + d.x + "," + d.y + ")";
            });
        });
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
                console.log(data);
                drawGraph(data);
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
