import { useState, useEffect, useRef } from "react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { Label } from "./ui/label";
import { Select, SelectItem } from "./ui/select";
import { toast } from "sonner";
import { Alert, AlertDescription, AlertTitle } from "./ui/alert";
import { Info, Terminal, AlertCircle, CheckCircle2, Rocket } from "lucide-react";

interface EnvVar {
  key: string;
  value: string;
}

interface AppTypePreset {
  installCmd: string;
  buildCmd: string;
  startCmd: string;
  port: number;
  notes?: string;
}

const APP_PRESETS: Record<string, AppTypePreset> = {
  node: {
    installCmd: "npm install",
    buildCmd: "",
    startCmd: "npm start",
    port: 3000,
    notes: "Standard Node.js application with npm"
  },
  react: {
    installCmd: "npm install",
    buildCmd: "npm run build",
    startCmd: "npm start",
    port: 3000,
    notes: "Create React App setup"
  },
  nextjs: {
    installCmd: "npm install",
    buildCmd: "npm run build",
    startCmd: "npm start",
    port: 3000,
    notes: "Next.js application (ensure output is 'standalone' for optimal deployment)"
  },
  flask: {
    installCmd: "pip install -r requirements.txt",
    buildCmd: "",
    startCmd: "python app.py",
    port: 5000,
    notes: "Flask applications require a WSGI server like Gunicorn for production"
  },
  go: {
    installCmd: "go mod download",
    buildCmd: "go build",
    startCmd: "./app",
    port: 8080,
    notes: "Go applications should build a single binary"
  },
  java: {
    installCmd: "mvn install",
    buildCmd: "mvn package",
    startCmd: "java -jar target/*.jar",
    port: 8080,
    notes: "Standard Spring Boot application"
  },
    vite: {
    installCmd: "npm install",
    buildCmd: "npm run build",
    startCmd: "npx serve dist",
    port: 3000,
    notes: "Vite applications are served using the preview command after build"
    }
};

export default function DeployForm() {
  const [repoUrl, setRepoUrl] = useState("");
  const logContainerRef = useRef<HTMLDivElement>(null);
  const [appType, setAppType] = useState("node");
  const [installCmd, setInstallCmd] = useState("");
  const [buildCmd, setBuildCmd] = useState("");
  const [startCmd, setStartCmd] = useState("");
  const [port, setPort] = useState(3000);
  const [rootFolder, setRootFolder] = useState("");
  const [envVars, setEnvVars] = useState<EnvVar[]>([{ key: "", value: "" }]);
  const [loading, setLoading] = useState(false);
  const [logs, setLogs] = useState<string[]>([]);
  const [deploymentStarted, setDeploymentStarted] = useState(false);
  const [deployedUrl, setDeployedUrl] = useState("");

  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const preset = APP_PRESETS[appType];
    if (preset) {
      setInstallCmd(preset.installCmd);
      setBuildCmd(preset.buildCmd);
      setStartCmd(preset.startCmd);
      setPort(preset.port);
    }
  }, [appType]);


  

  const handleEnvChange = (index: number, field: keyof EnvVar, value: string) => {
    const updated = [...envVars];
    updated[index][field] = value;
    setEnvVars(updated);
  };

  const addEnvVar = () => {
    setEnvVars([...envVars, { key: "", value: "" }]);
  };

  const removeEnvVar = (index: number) => {
    setEnvVars((prev) => prev.filter((_, i) => i !== index));
  };

  const validateForm = () => {
    if (!repoUrl) {
      toast.error("Repository URL is required");
      return false;
    }
    
    try {
      new URL(repoUrl);
    } catch (e) {
      toast.error("Please enter a valid URL");
      return false;
    }

    if (!startCmd) {
      toast.error("Start command is required");
      return false;
    }

    return true;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;

    const env: Record<string, string> = {};
    envVars.forEach(({ key, value }) => {
      if (key) env[key] = value;
    });

    const payload = {
      repoUrl,
      appType,
      installCmd,
      buildCmd,
      startCmd,
      port: Number(port),
      rootFolder,
      env,
    };

    try {
      setLogs([]);
      setLoading(true);
      setDeploymentStarted(true);
      setError(null);
      setDeployedUrl("");

      const streamRes = await fetch("http://localhost:8080/build-stream", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!streamRes.ok) {
        throw new Error(`Server responded with ${streamRes.status}`);
      }

      if (!streamRes.body) {
        throw new Error("No response body received");
      }

      const reader = streamRes.body.getReader();
      const decoder = new TextDecoder("utf-8");

      while (true) {
        const { value, done } = await reader.read();
        if (done) break;
        
        const chunk = decoder.decode(value);
        const lines = chunk.split("\n\n").filter(Boolean);

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const log = line.replace("data: ", "");
            setLogs((prev) => [...prev, log]);

          
            const urlMatch = log.match(/http?:\/\/localhost:[0-9]+[^\s]*/);
            if (urlMatch) {
              setDeployedUrl(urlMatch[0]);

              toast.success("Deployment Successful", {
                description: `Your application is now live at ${urlMatch[0]}`,
                action: {
                  label: "Open",
                  onClick: () => window.open(urlMatch[0], "_blank")
                }
              });
            }

           
            if (log.toLowerCase().includes("error") || log.toLowerCase().includes("failed")) {
              setError(log);
            }
          }
        }
      }
    } catch (err: any) {
      console.error("Deployment error:", err);
      const errorMessage = err?.message || "Deployment failed due to an unknown error";
      setError(errorMessage);
      toast.error("Deployment Failed", { 
        description: errorMessage,
        duration: 10000 
      });
    } finally {
      setLoading(false);
      
    }
  };

  useEffect(() => {
    const el = logContainerRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [logs]);

  const resetForm = () => {
    setDeploymentStarted(false);
    setDeployedUrl("");
    setLogs([]);
    setError(null);
  };

  const getCurrentPresetNotes = () => {
    return APP_PRESETS[appType]?.notes || "No specific notes for this application type.";
  };

  

  return (
    <div className="min-h-screen bg-background text-foreground">
      <style jsx>{`
        ::placeholder {
          color: oklch(0.45 0.02 240) !important;
        }
        input::placeholder {
          color: oklch(0.45 0.02 240) !important;
        }
        select option {
          background-color: oklch(0.16 0.02 240);
          color: oklch(0.94 0.01 240);
        }
        .bg-background {
          background-color: oklch(0.08 0.01 240);
        }
        .text-foreground {
          color: oklch(0.94 0.01 240);
        }
        .bg-card {
          background-color: oklch(0.12 0.015 240);
        }
        .border {
          border-color: oklch(0.22 0.03 240);
        }
        .text-muted-foreground {
          color: oklch(0.65 0.02 240);
        }
        .bg-muted {
          background-color: oklch(0.18 0.02 240);
        }
        .bg-input {
          background-color: oklch(0.16 0.02 240);
        }
        .bg-primary {
          background-color: oklch(0.65 0.25 260);
        }
        .text-primary {
          color: oklch(0.65 0.25 260);
        }
        .bg-destructive {
          background-color: oklch(0.62 0.25 20);
        }
        .bg-secondary {
          background-color: oklch(0.20 0.02 240);
        }
        .text-secondary-foreground {
          color: oklch(0.92 0.01 240);
        }
        .border-primary {
          border-color: oklch(0.65 0.25 260);
        }
        .hover\\:bg-muted:hover {
          background-color: oklch(0.18 0.02 240);
        }
        .hover\\:text-primary\\/80:hover {
          color: oklch(0.65 0.25 260 / 0.8);
        }
        .bg-green-500\\/10 {
          background-color: oklch(0.70 0.15 120 / 0.1);
        }
        .border-green-500\\/30 {
          border-color: oklch(0.70 0.15 120 / 0.3);
        }
        .text-green-400 {
          color: oklch(0.70 0.15 120);
        }
        .hover\\:border-primary\\/50:hover {
          border-color: oklch(0.65 0.25 260 / 0.5);
        }
        .focus\\:ring-primary:focus {
          --tw-ring-color: oklch(0.65 0.25 260);
        }
        .focus\\:border-primary:focus {
          border-color: oklch(0.65 0.25 260);
        }
      `}</style>
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
        
          <div className="text-center mb-12">
            <h1 className="text-5xl font-bold bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent mb-4">
              Deploy Your Application
            </h1>
            <p className="text-muted-foreground text-xl max-w-2xl mx-auto">
              Deploy your code with ease using our streamlined platform. Fast, reliable, and secure.
            </p>
          </div>

          <div className="flex flex-col gap-8 justify-around items-center">
           
            <div className="lg:w-3/4 space-y-6">
              <div className="rounded-xl border bg-card shadow-lg p-6">
                {!deploymentStarted ? (
                  <div className="space-y-6">
                    <div>
                      <Alert className="mb-6">
                        <Terminal className="h-4 w-4" />
                        <AlertTitle>Quick Tip</AlertTitle>
                        <AlertDescription>
                          {getCurrentPresetNotes()}
                        </AlertDescription>
                      </Alert>

                      {error && (
                        <Alert variant="destructive" className="mb-6">
                          <AlertCircle className="h-4 w-4" />
                          <AlertTitle>Error</AlertTitle>
                          <AlertDescription>
                            {error}
                          </AlertDescription>
                        </Alert>
                      )}
                    </div>

                    <div className="space-y-3">
                      <Label className="text-base font-semibold flex items-center gap-2">
                        <Info className="h-4 w-4" />
                        Git Repository URL
                      </Label>
                      <Input
                        value={repoUrl}
                        onChange={(e) => setRepoUrl(e.target.value)}
                        placeholder="https://github.com/username/repository.git"
                        className="h-12 text-base bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                      />
                    </div>

                    <div className="space-y-3">
                      <Label className="text-base font-semibold">Application Type</Label>
                      <Select 
                        className="h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary" 
                        value={appType} 
                        onValueChange={setAppType}
                      >
                        <SelectItem value="node">Node.js</SelectItem>
                        <SelectItem value="react">React</SelectItem>
                        <SelectItem value="nextjs">Next.js</SelectItem>
                        <SelectItem value="flask">Flask</SelectItem>
                        <SelectItem value="go">Go</SelectItem>
                        <SelectItem value="java">Java</SelectItem>
                        <SelectItem value="vite">Vite</SelectItem>
                      </Select>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-3">
                        <Label className="text-base font-semibold">Install Command</Label>
                        <Input
                          value={installCmd}
                          onChange={(e) => setInstallCmd(e.target.value)}
                          placeholder="npm install"
                          className="h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                        />
                      </div>
                      <div className="space-y-3">
                        <Label className="text-base font-semibold">Build Command</Label>
                        <Input
                          value={buildCmd}
                          onChange={(e) => setBuildCmd(e.target.value)}
                          placeholder="npm run build"
                          className="h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-3">
                        <Label className="text-base font-semibold">Start Command</Label>
                        <Input
                          value={startCmd}
                          onChange={(e) => setStartCmd(e.target.value)}
                          placeholder="npm start"
                          className="h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                        />
                      </div>
                      <div className="space-y-3">
                        <Label className="text-base font-semibold">Port</Label>
                        <Input
                          type="number"
                          value={port}
                          onChange={(e) => setPort(Number(e.target.value))}
                          placeholder="3000"
                          className="h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                        />
                      </div>
                    </div>

                    <div className="space-y-3">
                      <Label className="text-base font-semibold">Root Folder (optional)</Label>
                      <Input
                        value={rootFolder}
                        onChange={(e) => setRootFolder(e.target.value)}
                        placeholder="Leave empty if root directory"
                        className="h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                      />
                    </div>

                    <div className="space-y-4">
                      <Label className="text-base font-semibold">Environment Variables</Label>
                      <div className="space-y-3">
                        {envVars.map((pair, index) => (
                          <div key={index} className="flex gap-3 items-center">
                            <Input
                              placeholder="KEY"
                              value={pair.key}
                              onChange={(e) => handleEnvChange(index, "key", e.target.value)}
                              className="flex-1 h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                            />
                            <Input
                              placeholder="VALUE"
                              value={pair.value}
                              onChange={(e) => handleEnvChange(index, "value", e.target.value)}
                              className="flex-1 h-12 bg-input border focus:border-primary focus:ring-1 focus:ring-primary"
                            />
                            {envVars.length > 1 && (
                              <Button
                                variant="destructive"
                                size="icon"
                                onClick={() => removeEnvVar(index)}
                                className="h-12 w-12 shrink-0"
                              >
                                ×
                              </Button>
                            )}
                          </div>
                        ))}
                      </div>
                      <Button
                        variant="outline"
                        onClick={addEnvVar}
                        className="w-full h-12 text-base font-semibold border hover:bg-muted"
                      >
                        + Add Environment Variable
                      </Button>
                    </div>

                    <Button
                      onClick={handleSubmit}
                      disabled={loading || !repoUrl}
                      className="w-full h-14 text-lg font-bold bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 shadow-lg hover:shadow-xl transition-all duration-200"
                    >
                      {loading ? (
                        <div className="flex items-center justify-center">
                          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-white mr-3"></div>
                          Deploying...
                        </div>
                      ) : (
                        <span className="flex items-center">
                          <Rocket className="mr-3 h-5 w-5" />
                          Deploy Application
                        </span>
                      )}
                    </Button>
                  </div>
                ) : (
                  <div className="text-center py-12">
                    {!deployedUrl && (
                      <>
                        <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-primary mx-auto mb-6"></div>
                        <h3 className="text-2xl font-bold mb-3">Deployment in Progress</h3>
                        <p className="text-muted-foreground text-lg mb-6">
                          Please wait while we deploy your application...
                        </p>
                      </>
                    )}
                    
                    {deployedUrl && (
                      <div className="mt-6 p-6 bg-green-500/10 rounded-xl border border-green-500/30">
                        <CheckCircle2 className="h-8 w-8 text-green-400 mx-auto mb-3" />
                        <p className="text-green-400 mb-3 text-lg font-semibold">
                          Deployment Successful!
                        </p>
                        <a
                          href={deployedUrl}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-primary hover:text-primary/80 underline text-lg font-medium"
                        >
                          {deployedUrl}
                        </a>
                      </div>
                    )}

                    <Button 
                      onClick={resetForm} 
                      variant="outline" 
                      className="mt-6"
                    >
                      Start New Deployment
                    </Button>
                  </div>
                )}
              </div>
            </div>
        
            
            <div className="lg:w-3/4 w-full">
              <div className="rounded-xl border bg-card shadow-lg h-full">
                <div className="p-4 border-b">
                  <h3 className="font-semibold text-lg">Deployment Logs</h3>
                </div>
                
                <div className="h-[600px] overflow-auto" ref={logContainerRef}>
                  {logs.length > 0 ? (
                    <div className="p-4">
                      {logs.map((log, index) => (
                        <div key={index} className="text-sm text-muted-foreground font-mono mb-1">
                          {log}
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="p-12 text-center text-muted-foreground">
                      <div className="text-6xl mb-4">📋</div>
                      <p className="text-lg">Deployment logs will appear here...</p>
                      {deploymentStarted && (
                        <p className="mt-2 text-sm">Waiting for logs from the deployment server...</p>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}