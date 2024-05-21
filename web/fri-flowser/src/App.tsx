import './App.css'
import Header from './components/Header';
import configureCadence from './lib/candance'; 

import { useEffect, useRef } from 'react'
import Editor, { DiffEditor, useMonaco, loader } from '@monaco-editor/react';

function App() {
  const LANGUAGE_CADENCE = 'cadence';
  const editorRef = useRef(null);

  const beforeEditorMount = (monaco) => {
    configureCadence(monaco);
  }

  const handleEditorDidMount = (editor) => {

    editorRef.current = editor;


    console.log('Editor instance:', editor);

  };

  const setEditorContent = (content) => {
    if (editorRef.current) {
      const model = editorRef.current.getModel(); // Ensure getModel method exists
      if (model) {
        model.setValue(content);
        console.log(content + ' test');
      } else {
        console.error('Model is undefined');
      }
    } else {
      console.error('Editor instance is undefined');
    }
  };

  return (
    <div className='flex flex-col'>
      <Header setEditorContent={setEditorContent} />

      <Editor
        theme='vs-dark'
        language={LANGUAGE_CADENCE}
        value={"//candance code"}
        style="pt-20"
        height="95vh"
        options={{
          automaticLayout: true
        }}
        onMount={handleEditorDidMount}
        beforeMount={beforeEditorMount}
      />
    </div>
  )
}

export default App;