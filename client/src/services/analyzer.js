import axios from "axios";

const BASIC_ENDPOINT = 'http://localhost:8080/tokens';

export async function analyzeCode(code) {
    try {
        const response = await axios.post(BASIC_ENDPOINT, { query: code });
        return response.data;
    } catch (error) {
        console.error('Error analyzing code:', error);
        throw error;
    }
}