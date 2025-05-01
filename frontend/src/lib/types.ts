export interface KeyShare {
    group: string;
    sk: string;
    public: string;
    details: {
        secret: string;
        verificationKey: string;
        publicKey: string;
        vssCommitment: string[];
        id: number;
        group: number;
    };
}

export interface KeyData {
    privateKey: string;
    keyShares: KeyShare[];
    verificationKeys: string[];
    savedAt: string;
}

export interface CommitKeyShareResponse {
    signer: string;
    commitment: string;
} 