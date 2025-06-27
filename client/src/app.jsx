import './app.css';
import { useState } from 'react';
import { analyzeCode } from './services/analyzer';
import Token from './components/token/token';

export default function App() {
    const [code, setCode] = useState('');
    const [results, setResults] = useState([]);
    const handleChange = (event) => { setCode(event.target.value); };

    const handleAnalyzer = async () => {
        try {
            const response = await analyzeCode(code);
            console.log(response)
            setResults(response);
        } catch (error) {
            console.error('Error al analizar el código:', error);
        }
    }
    const handleClear = () => {
        setCode('');
        const textarea = document.querySelector('.codeSection');
        if (textarea) textarea.value = '';

    };

    const handleKeyDown = (event) => {
        if (event.key === 'Tab') {
            event.preventDefault();
            const textarea = event.target;
            const start = textarea.selectionStart;
            const end = textarea.selectionEnd;

            const newValue = code.substring(0, start) + '\t' + code.substring(end);
            setCode(newValue);

            requestAnimationFrame(() => {
                textarea.selectionStart = textarea.selectionEnd = start + 1;
            });
        }
    };



    return (
        <div className="App page">
            <textarea onChange={handleChange} className='codeSection' onKeyDown={handleKeyDown}></textarea>
            <div className='controllers'>
                <button className='an' onClick={handleAnalyzer}>Analizar</button>
                <button className='cl' onClick={handleClear}>Limpiar</button>
            </div>

            <div className='results'>
                <div className='Conteo Cabeza'>
                    <p>ID</p>
                    <p>PR</p>
                    <p>Simbolo</p>
                    <p>Numero</p>
                    <p>Cadena</p>
                </div>
                {
                    results && (
                        <div className='Cabeza'>
                            <p>{results?.conteo?.ID}</p>
                            <p>{results?.conteo?.PR}</p>
                            <p>{results?.conteo?.Símbolo}</p>
                            <p>{results?.conteo?.Número}</p>
                            <p>{results?.conteo?.Cadenas}</p>
                        </div>
                    )
                }
                <div className='Cabeza'>
                    <p>Token</p>
                    <p>Tipo</p>
                    <p>Linea</p>
                </div>
                {
                    results && (
                        results?.tokens?.map((token, index) => (
                            <Token key={index} data={token} />
                        ))
                    )
                }
            </div>

            {
                results && (
                    results?.semantic?.map((val, index) => (
                        <div key={index} className='semantic'>
                            {val}
                        </div>
                    ))
                )
            }
            {
                results && (
                    results?.syntax?.map((val, index) => (
                        <div key={index} className='semantic'>
                            {val}
                        </div>
                    ))
                )
            }
        </div>
    );
}