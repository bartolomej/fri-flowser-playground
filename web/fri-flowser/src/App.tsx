import './App.css'
import Header from './components/Header';
import MonacoEditor from '@monaco-editor/react';
import configureCadence from './lib/candance'; 

function App() {
  const LANGUAGE_CADENCE = 'cadence';
  
  const handleEditorWillMount = (monaco) => {
    configureCadence(monaco);
  };


  return (
    <div className='flex flex-col'>
      <Header />
      <MonacoEditor
        theme='vs-dark'
        language={LANGUAGE_CADENCE}
        value={"//candance code"}
        style="pt-20"
        height="95vh"
        options={{
          automaticLayout: true
        }}
        beforeMount={handleEditorWillMount}
      />
    </div>
  )
}

export default App;