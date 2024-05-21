import './App.css'
import configureCadence from '@/common/candance';

import {useEffect, useState} from 'react'
import Editor, {Monaco} from '@monaco-editor/react';
import {ProjectFile, ProjectService} from "@/project.service.ts";

function App() {
    const LANGUAGE_CADENCE = 'cadence';
    const [code, setCode] = useState('');
    const [projectFiles, setProjectFiles] = useState<ProjectFile[]>();
    const service = new ProjectService({baseUrl: "http://localhost:8080"});
    const urlParams = new URLSearchParams(window.location.search);
    const projectUrl = urlParams.get('projectUrl');

    useEffect(() => {
        (async function () {
            if (projectUrl) {
                await service.openProject(projectUrl)
                setProjectFiles(await service.listProjectFiles())
            }
        })()
    }, [projectUrl]);

    const beforeEditorMount = (monaco: Monaco) => {
        configureCadence(monaco);
    }

    if (!projectUrl) {
        return (
            <div>
                Set `projectUrl` query parameter
            </div>
        )
    }

    return (
        <div className='flex flex-row'>
            <div className="flex flex-col justify-between p-2">
                <div className="flex flex-col gap-y-[10px]">
                    {projectFiles
                        ?.filter(file => !file.isDirectory)
                        .map(file => {
                            const fileName = file.path.split("/").reverse()[0]
                            return (
                                <div key={file.path} onClick={() => setCode(file.path)}
                                     className="max-w-[200px] truncate text-left">
                                    {fileName}
                                </div>
                            )
                        })}
                </div>

                {code && (
                    <div>
                        <label>
                            Arguments
                            <textarea rows={10}></textarea>
                        </label>
                        <button>Execute</button>
                    </div>
                )}
            </div>

            <Editor
                theme='vs-dark'
                language={LANGUAGE_CADENCE}
                value={code}
                onChange={code => setCode(code ?? "")}
                className="h-[100vh] pt-2 w-full"
                options={{automaticLayout: true}}
                beforeMount={beforeEditorMount}
            />
        </div>
    )
}

export default App;
