/*
import React from 'react';
import Tree from 'react-d3-tree';

function LineageChart({lineageData}) {
    return (
        <div id="treeWrapper">
            <Tree data={lineageData} />
        </div>
    )
}

export default LineageChart
*/

import Tree from 'react-tree-graph';

function LineageChart({lineageData, height, width}) {

    const handleClick = (event, node) => {
        window.location = `/asset/${node}`;
      }

    return (
        <div className="custom-container">
            <Tree 
            data={lineageData} 
            height={height} 
            width={width}
            nodeRadius={10}
            svgProps={{className: 'customlineagechart'}}
            gProps={{onClick: handleClick}}
            />
        </div>
    ) 
}

export default LineageChart
