"use client";

import { useEffect, useState } from "react";

interface NodeResponse {
  hostname: string;
  status: string;
  assigned_ssh_port: number;
  last_seen_at: string;
}

export default function Home() {
  const [nodes, setNodes] = useState<NodeResponse[]>([]);
  const [loading, setLoading] = useState(true);

  // Poll the Go Central Server for nodes
  useEffect(() => {
    const fetchNodes = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/v1/nodes");
        if (res.ok) {
          const data = await res.json();
          setNodes(data);
        }
      } catch (err) {
        console.error("Failed to fetch nodes", err);
      } finally {
        setLoading(false);
      }
    };

    fetchNodes();
    // Phase 4.1 uses simple 3s polling
    const interval = setInterval(fetchNodes, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-[#0a0a0a] text-white p-8 md:p-16 font-sans selection:bg-emerald-500/30">
      <header className="mb-16">
        <h1 className="text-5xl font-extrabold tracking-tight bg-gradient-to-r from-blue-400 to-emerald-400 bg-clip-text text-transparent drop-shadow-sm mb-4">
          ByoHIL Control Center
        </h1>
        <p className="text-gray-400 text-lg max-w-2xl leading-relaxed">
          Monitor your active geographically distributed Hardware-in-the-Loop test bench topologies in real-time.
        </p>
      </header>

      {loading ? (
        <div className="flex animate-pulse space-x-6">
          <div className="rounded-2xl bg-gray-800/50 border border-gray-800 h-48 w-full md:w-1/3 shadow-sm"></div>
          <div className="rounded-2xl bg-gray-800/50 border border-gray-800 h-48 w-full md:w-1/3 shadow-sm hidden md:block"></div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
          {nodes.length === 0 ? (
            <div className="col-span-full py-12 flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-700 bg-gray-900/20">
              <svg className="w-16 h-16 text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M5 12h14M12 5l7 7-7 7" />
              </svg>
              <p className="text-gray-400 text-lg font-medium">No bare-metal nodes found online.</p>
              <p className="text-gray-500 text-sm mt-2">Deploy the Python node agent to see your test benches register here!</p>
            </div>
          ) : (
            nodes.map((node) => (
              <div
                key={node.hostname}
                className="group relative overflow-hidden rounded-2xl border border-gray-800 bg-gray-900/60 p-7 shadow-2xl transition-all duration-300 hover:border-gray-600 hover:shadow-emerald-900/20"
              >
                {/* Status Indicator Bar */}
                <div 
                  className={`absolute top-0 left-0 w-full h-1 transition-all duration-500 ${
                    node.status === "online"
                      ? "bg-gradient-to-r from-emerald-400 to-emerald-600"
                      : node.status === "in-use"
                      ? "bg-gradient-to-r from-amber-400 to-amber-600"
                      : "bg-gradient-to-r from-red-500 to-red-600"
                  }`}
                />

                <div className="flex justify-between items-start mb-6 mt-2 relative z-10">
                  <h2 className="text-2xl font-bold tracking-tight text-gray-100 group-hover:text-white transition-colors">
                    {node.hostname}
                  </h2>
                  <span
                    className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wider shadow-sm backdrop-blur-sm ${
                      node.status === "online"
                        ? "bg-emerald-400/10 text-emerald-400 border border-emerald-400/20"
                        : node.status === "in-use"
                        ? "bg-amber-400/10 text-amber-400 border border-amber-400/20"
                        : "bg-red-400/10 text-red-400 border border-red-400/20"
                    }`}
                  >
                    <span className={`w-1.5 h-1.5 rounded-full mr-2 ${
                      node.status === "online" ? "bg-emerald-400 animate-[pulse_2s_ease-in-out_infinite]" : "bg-current"
                    }`} />
                    {node.status}
                  </span>
                </div>

                <div className="space-y-4 relative z-10">
                  <div className="flex items-center justify-between p-3 rounded-xl bg-black/50 border border-gray-800/80 shadow-inner">
                    <span className="text-gray-400 text-sm font-medium">Routing Port</span>
                    <span className="font-mono text-emerald-300 bg-emerald-950/50 px-3 py-1 rounded-lg text-sm font-bold tracking-widest border border-emerald-800/50 shadow-sm">
                      {node.assigned_ssh_port}
                    </span>
                  </div>
                  
                  <div className="flex items-center justify-between px-3 text-sm">
                    <span className="text-gray-500 font-medium">Last Keepalive</span>
                    <span className="text-gray-300 tabular-nums bg-gray-800/50 px-2 py-0.5 rounded">
                      {new Date(node.last_seen_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
                    </span>
                  </div>
                </div>

                {/* Subtle Ambient Glow */}
                <div
                  className={`absolute -bottom-20 -right-20 w-48 h-48 rounded-full blur-[80px] opacity-20 pointer-events-none transition-all duration-700 group-hover:opacity-40 group-hover:scale-110 ${
                    node.status === "online"
                      ? "bg-emerald-500"
                      : node.status === "in-use"
                      ? "bg-amber-500"
                      : "bg-red-500"
                  }`}
                />
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
