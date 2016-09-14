require(stringr)
require(igraph)
require(plyr)
require(reshape2)

rawData <- read.csv("data/wiki-mouse-2.csv")

rawData$page <- gsub("http://en.wikipedia.org/wiki/", "", rawData$page)
rawData$link <- gsub("http://en.wikipedia.org/wiki/", "", rawData$link)

power <- 14
size <- 2^power
version <- "v0.5"
fileName <- paste("renders/wiki-mouse-2-", power, "-", version, ".png",sep="")


edgeList <- simplify(graph.data.frame(rawData,directed=TRUE))
png(filename = fileName,
    width = size,
    height = size,
    bg="black")

plot.igraph(edgeList,
            #layout=layout.kamada.kawai,
            #layout = layout.fruchterman.reingold,
            vertex.size=0.1,
            edge.width=0.5,
            edge.color="white",
            vertex.label.color="red",
            vertex.label.cex = 1
            )
title(main=paste("Links -- Rendered at ",size,"^2.  ",length(rawData$page), version,sep=""), cex.main=7)

dev.off()
