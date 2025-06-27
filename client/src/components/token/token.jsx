import './token.css';

export default function Token({ data }) {
    return (
        <div className='Token'>
            <p>{data.token}</p>
            <p>{data.tipo}</p>
            <p>{data.linea}</p>
        </div>
    )
}