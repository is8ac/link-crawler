require(stringr)
require(igraph)
require(plyr)
require(reshape2)

rawData <- read.csv("data/ssc-all.csv")


power <- 13
size <- 2^power
version <- "v0.3"
fileName <- paste("renders/ssc-all-kk-", power, "-", version, ".png",sep="")


edgeList <- graph.data.frame(rawData,directed=TRUE)
png(filename = fileName,
    width = size,
    height = size,
    bg="black")

plot.igraph(edgeList,
            vertex.size=0.5,
            edge.width=1,
            edge.color="white",
            vertex.label.color="white",
            vertex.label.cex = 0.9
            )
title(main=paste("Links -- Rendered at ",size,"^2.  ",length(rawData$page), version,sep=""), cex.main=7)

dev.off()
