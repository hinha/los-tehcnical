
describe('Loan Lifecycle API Test', () => {
    let loanId = null;

    it('should create a loan', () => {
        cy.request('POST', '/loans', {
            borrower_id: 'B001',
            principal_amount: 5000000,
            rate: 10,
            roi: 15
        }).then((response) => {
            expect(response.status).to.eq(201);
            expect(response.body.data).to.have.property('id');
            loanId = response.body.data.id;
        });
    });

    it('should not allow duplicate borrower ID', () => {
        cy.request({
            method: 'POST',
            url: '/loans',
            body: {
                borrower_id: 'B001',
                principal_amount: 6000000,
                rate: 12,
                roi: 18
            },
            failOnStatusCode: false
        }).then((response) => {
            expect(response.body).to.have.property('message', 'Failed to create loan');
            expect(response.body).to.have.property('errors');
            expect(response.body.errors).to.include('already exists');
        });
    });

    it('should approve the loan', () => {
        cy.request('POST', `/loans/${loanId}/approve`, {
            validator_id: 'VALID-001',
            proof_url: 'https://storage.yourstorage.com/loan-proof/visit123.jpeg'
        }).then((response) => {
            expect(response.status).to.eq(200);
        });
    });

    it('should add an investment', () => {
        cy.request('POST', `/loans/${loanId}/invest`, {
            investor_id: 'INV-001',
            email: 'investor@example.com',
            amount: 5000000
        }).then((response) => {
            expect(response.status).to.eq(200);
        });
    });

    it('should disburse the loan', () => {
        cy.request('POST', `/loans/${loanId}/disburse`, {
            field_officer_id: 'FO-001',
            signed_agreement: 'https://storage.yourstorage.com/loan-aggreement/aggreement.pdf'
        }).then((response) => {
            expect(response.status).to.eq(200);
        });
    });
});