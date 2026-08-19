import { ReactFlowProvider } from "@xyflow/react";
import { Canvas } from "./Canvas";
import "./layout.css";

function App() {
  return (
    <ReactFlowProvider>
      <Canvas />
    </ReactFlowProvider>
  );
}

export default App;
