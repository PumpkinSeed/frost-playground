export interface StoredCommitData {
    secret: string;
    signer: string;
    commitment: string;
    savedAt: string;
}

export function findCommitDataBySecret(secret: string): StoredCommitData | undefined {
    try {
        const commitData: StoredCommitData[] = JSON.parse(localStorage.getItem('commit_data') || '[]');
        return commitData.find(data => data.secret === secret);
    } catch (err) {
        console.error('Error finding commit data:', err);
        return undefined;
    }
}

export function getAllCommitData(): StoredCommitData[] {
    try {
        return JSON.parse(localStorage.getItem('commit_data') || '[]');
    } catch (err) {
        console.error('Error getting commit data:', err);
        return [];
    }
}

export function clearCommitData(): void {
    try {
        localStorage.removeItem('commit_data');
    } catch (err) {
        console.error('Error clearing commit data:', err);
    }
} 