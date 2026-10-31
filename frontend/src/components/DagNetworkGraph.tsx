import { useMemo, useCallback, memo } from 'react';
import { ReactFlow, Background, Controls, MarkerType, Position, Handle, ReactFlowProvider } from '@xyflow/react';
import type { Node, Edge as FlowEdge, NodeProps } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { Box, Typography, Button } from '@mui/material';
import DownloadIcon from '@mui/icons-material/Download';
import dagre from 'dagre';
import { toPng } from 'html-to-image';

interface Operation {
  id: number;
  name: string;
  machine_type: number;
  processing_time: string;
}

interface Edge {
  from_op_id: number;
  to_op_id: number;
}

interface DagNetworkGraphProps {
  operations: Operation[];
  edges: Edge[];
  machineColorMap: Record<number, string>;
}

// 1. Створюємо справжній кастомний вузол React Flow.
// Тепер бібліотека знає його точні розміри (ширина 240px, висота ~110px), і нічого не накладається.
const MachineNode = memo(({ data }: NodeProps) => {
  const { op, machineColor } = data as { op: Operation; machineColor: string };

  return (
    <Box
      sx={{
        p: '14px 16px',
        borderRadius: '8px',
        backgroundColor: '#fff',
        border: '1px solid #dadce0',
        borderLeft: `5px solid ${machineColor}`,
        boxShadow: '0 4px 12px rgba(0,0,0,0.08)',
        width: '240px',
        position: 'relative',
      }}
    >
      {/* Лівий порт (вхід) */}
      <Handle
        type="target"
        position={Position.Left}
        style={{
          background: '#9aa0a6',
          width: '8px',
          height: '8px',
          left: '-4px',
          top: '50%',
          transform: 'translateY(-50%)',
          border: '2px solid #fff',
          borderRadius: '50%',
        }}
      />
      {/* Правий порт (вихід) */}
      <Handle
        type="source"
        position={Position.Right}
        style={{
          background: '#9aa0a6',
          width: '8px',
          height: '8px',
          right: '-4px',
          top: '50%',
          transform: 'translateY(-50%)',
          border: '2px solid #fff',
          borderRadius: '50%',
        }}
      />

      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: '4px' }}>
        <Typography sx={{ fontSize: '11px', fontWeight: 700, color: '#5f6368' }}>
          Op #{op.id}
        </Typography>
        <Typography
          sx={{
            fontSize: '10px',
            fontWeight: 600,
            color: machineColor,
            backgroundColor: `${machineColor}15`,
            padding: '2px 6px',
            borderRadius: '4px',
          }}
        >
          Type {op.machine_type}
        </Typography>
      </Box>
      <Typography sx={{ fontSize: '13px', fontWeight: 600, color: '#202124', mb: '6px', lineHeight: 1.3 }}>
        {op.name}
      </Typography>
      <Typography sx={{ fontSize: '12px', color: '#2e7d32', fontWeight: 500 }}>
        ⏱ {op.processing_time}
      </Typography>
    </Box>
  );
});

MachineNode.displayName = 'MachineNode';

// Реєструємо типи вузлів поза компонентом, щоб уникнути зайвих ререндерів
const nodeTypes = {
  machineNode: MachineNode,
};

// Авторозміщення через Dagre для горизонтального графа (Left to Right)
const getLayoutedElements = (nodes: Node[], edges: FlowEdge[]) => {
  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));

  dagreGraph.setGraph({ rankdir: 'LR', nodesep: 70, ranksep: 140 });

  const nodeWidth = 240;
  const nodeHeight = 110;

  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  const layoutedNodes = nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    return {
      ...node,
      position: {
        x: nodeWithPosition.x - nodeWidth / 2,
        y: nodeWithPosition.y - nodeHeight / 2,
      },
    };
  });

  return { nodes: layoutedNodes, edges };
};

// Внутрішній компонент для доступу до хуків React Flow
function FlowContent({ operations, edges, machineColorMap }: DagNetworkGraphProps) {
  const { nodes, edges: flowEdges } = useMemo(() => {
    const initialNodes: Node[] = [];
    const initialEdges: FlowEdge[] = [];

    operations.forEach((op) => {
      const machineColor = machineColorMap[op.machine_type] || '#5f6368';

      initialNodes.push({
        id: String(op.id),
        type: 'machineNode', // Використовуємо наш власний тип кастомного вузла
        position: { x: 0, y: 0 },
        data: { op, machineColor }, // Передаємо дані чистим об'єктом
      });
    });

    edges.forEach((edge, idx) => {
      initialEdges.push({
        id: `e-${idx}`,
        source: String(edge.from_op_id),
        target: String(edge.to_op_id),
        type: 'default', // Плавні криві Безьє
        animated: false,
        style: { stroke: '#1a73e8', strokeWidth: 2.5 },
        markerEnd: {
          type: MarkerType.ArrowClosed,
          color: '#1a73e8',
          width: 20,
          height: 20,
        },
      });
    });

    return getLayoutedElements(initialNodes, initialEdges);
  }, [operations, edges, machineColorMap]);

  // Функція збереження графа у PNG
  const handleDownloadImage = useCallback(() => {
    const viewportElement = document.querySelector('.react-flow__viewport') as HTMLElement;
    if (!viewportElement) return;

    toPng(viewportElement, {
      cacheBust: true,
      backgroundColor: '#f8f9fa',
    })
      .then((dataUrl) => {
        const link = document.createElement('a');
        link.download = 'dag-topology-graph.png';
        link.href = dataUrl;
        link.click();
      })
      .catch((err) => {
        console.error('Failed to generate image', err);
      });
  }, []);

  return (
    <Box sx={{ width: '100%', height: '550px', border: '1px solid #dadce0', borderRadius: '8px', background: '#f8f9fa', position: 'relative' }}>

      {/* Кнопка завантаження картинки */}
      <Box sx={{ position: 'absolute', top: '16px', right: '16px', zIndex: 10 }}>
        <Button
          variant="outlined"
          size="small"
          onClick={handleDownloadInspector => handleDownloadImage()}
          startIcon={<DownloadIcon />}
          sx={{
            backgroundColor: '#fff',
            borderColor: '#dadce0',
            color: '#3c4043',
            textTransform: 'none',
            fontSize: '12px',
            fontWeight: 500,
            boxShadow: '0 1px 2px rgba(0,0,0,0.05)',
            '&:hover': {
              backgroundColor: '#f1f3f4',
              borderColor: '#1a73e8',
              color: '#1a73e8',
            },
          }}
        >
          Export as PNG
        </Button>
      </Box>

      <ReactFlow
        nodes={nodes}
        edges={flowEdges}
        nodeTypes={nodeTypes} // Передаємо наш реєстр кастомних вузлів
        fitView
        attributionPosition="bottom-left"
        minZoom={0.2}
      >
        <Background gap={20} size={1} color="#e0e0e0" />
        <Controls />
      </ReactFlow>
    </Box>
  );
}

export function DagNetworkGraph(props: DagNetworkGraphProps) {
  return (
    <ReactFlowProvider>
      <FlowContent {...props} />
    </ReactFlowProvider>
  );
}
