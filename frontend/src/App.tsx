import { Toaster } from "sonner";

import DeployForm from "./components/Form";

function App() {
  return (
    <>
      <DeployForm />
      <Toaster richColors position="top-right" />
    </>
  );
}

export default App;