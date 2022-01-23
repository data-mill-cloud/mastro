import React from 'react';
import ReactFlow from 'react-flow-renderer';

//const onLoad = (reactFlowInstance) => reactFlowInstance.fitView();
const onElementClick = (event, element) => window.location = `/asset/${element.data.label}`;

function LineageChart({lineageData}) {

    return (
            <ReactFlow 
            elements={lineageData} 
            onElementClick={onElementClick}
            zoomOnScroll={true}
            //onLoad={onLoad}
            selectNodesOnDrag={false}
            defaultZoom={1}
            > </ReactFlow>      
    );
  };
  
  export default LineageChart;