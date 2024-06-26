function generateGUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        let r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// Action Submit Specialist Make
function submitReadSpecialistForm(formId) {
    this.event.preventDefault();
    const form = document.getElementById(formId);
    const formData = new FormData(form);
    const data = new FormData();

    formData.forEach((value, key) => {
        data.append(key, value)
    });

    let guid = generateGUID();
    fetch(form.action, {
        method: form.method, headers: {
            'x-correlation-id': guid
        }, body: data,
    })
        .then(response => {
                if (response.status === 500) {
                    throw new Error('Internal Server Error');
                }
                return response.json();
            }
        )
        .then(result => {
            if (result.status === 'Success') {
                document.getElementById('specialistProcess').style.display = 'block';
                document.getElementById('ExportSpecialist').style.display = 'block';
                alert('Successfully');
            } else {
                alert("Something went wrong !!! " + result.error)
            }
        })
        .catch(error => {
            console.error('Error:', error);
            if (error.message === 'Internal Server Error') {
                alert('An internal server error occurred. Please try again later.');
            } else {
                alert('An error occurred. Please try again.');
            }
        });
}


function submitQuerySpecialistForm(formId) {
    this.event.preventDefault();
    const form = document.getElementById(formId);
    const formData = new FormData(form);

    let guid = generateGUID();
    let url = form.action;

    const params = new URLSearchParams();
    formData.forEach((value, key) => {
        params.append(key, value);
    });
    url += '?' + params.toString();

    fetch(url, {
        method: 'GET', headers: {
            'x-correlation-id': guid
        }
    })
        .then(response => response.json())
        .then(result => {

            if (result.RowNo !== '') {
                document.getElementById('specialistDataTable').style.display = 'block';
                document.getElementById('RowNo').value = result.RowNo;
                document.getElementById('CustomerNo').value = result.CustomerNo;
                document.getElementById('CustomerNameTh').value = result.CustomerNameTh;
                document.getElementById('ChkThaiName').value = result.ChkThaiName;
                document.getElementById('ChkEngName').value = result.ChkEngName;
                document.getElementById('CustomerNameEn').value = result.CustomerNameEn;
                document.getElementById('Dob').value = result.Dob;
                document.getElementById('CountryCode').value = result.CountryCode;
                document.getElementById('CustomerType').value = result.CustomerType;
                document.getElementById('Zipcode').value = result.Zipcode;
                document.getElementById('Action').value = result.Action;
                document.getElementById('OldAct').value = result.OldAct;
                document.getElementById('Ovract').value = result.Ovract;
                document.getElementById('Pob').value = result.Pob;
                document.getElementById('ReasonCode').value = result.ReasonCode;
                document.getElementById('RtnCustomer').value = result.RtnCustomer;
                document.getElementById('SrcSeq').value = result.SrcSeq;
            } else {
                alert("Data not found")
                document.getElementById('specialistDataTable').style.display = 'none';
            }


        })
        .catch(error => {
            console.error('Error:', error);
        });
}


function resetForm() {
    document.getElementById('customerNoFilter').value = '';
    document.getElementById('RowNo').value = '';
    document.getElementById('CustomerNo').value = '';
    document.getElementById('CustomerNameTh').value = '';
    document.getElementById('ChkThaiName').value = '';
    document.getElementById('ChkEngName').value = '';
    document.getElementById('CustomerNameEn').value = '';
    document.getElementById('Dob').value = '';
    document.getElementById('CountryCode').value = '';
    document.getElementById('CustomerType').value = '';
    document.getElementById('Zipcode').value = '';
    document.getElementById('Action').value = '';
    document.getElementById('OldAct').value = '';
    document.getElementById('Ovract').value = '';
    document.getElementById('Pob').value = '';
    document.getElementById('ReasonCode').value = '';
    document.getElementById('RtnCustomer').value = '';
    document.getElementById('SrcSeq').value = '';
}

function defaultNewSpecialistForm() {
    document.getElementById('RowNo').value = '';
    document.getElementById('CustomerNo').value = '';
    document.getElementById('CustomerNameTh').value = '';
    document.getElementById('ChkThaiName').value = '';
    document.getElementById('ChkEngName').value = '';
    document.getElementById('CustomerNameEn').value = '';
    document.getElementById('Dob').value = '';
    document.getElementById('CountryCode').value = 'TH';
    document.getElementById('CustomerType').value = 'ZZ';
    document.getElementById('Zipcode').value = '0';
    document.getElementById('Action').value = 'R';
    document.getElementById('OldAct').value = 'B';
    document.getElementById('Ovract').value = 'Y';
    document.getElementById('Pob').value = '';
    document.getElementById('ReasonCode').value = '';
    document.getElementById('RtnCustomer').value = '';
    document.getElementById('SrcSeq').value = '(1)(2)(3)';
}

function newSpecialist() {
    document.getElementById('specialistDataTable').style.display = 'block';
    defaultNewSpecialistForm()
}

function cancelSpecialist() {
    document.getElementById('specialistDataTable').style.display = 'none';
    resetForm()
}

function saveSpecialist() {
    // Prevent default form submission
    this.event.preventDefault();

    // Create request body
    let requestBody = getDataFromInput();

    // Make the POST request
    fetch('/api/v1/specialist/save', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
    })
        .then(response => {
            if (response.status === 500) {
                throw new Error('Internal Server Error');
            }
            return response.json();
        })
        .then(result => {
            if (result.status === 'Success') {
                document.getElementById('specialistProcess').style.display = 'block';
                alert('Successfully');
                cancelSpecialist()
            } else {
                alert("Something went wrong !!! " + result.error)
            }
        })
        .catch((error) => {
            console.error('Error:', error);
            if (error.message === 'Internal Server Error') {
                alert('An internal server error occurred. Please try again later.');
            } else {
                alert('An error occurred. Please try again.');
            }
        });
}

function deleteSpecialist() {
    // Prevent default form submission
    this.event.preventDefault();
    const RowNo = document.getElementById('RowNo').value;
    let url = '/api/v1/specialist/delete' + '?rowNo=' + RowNo;
    // Make the POST request
    fetch(url, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            if (response.status === 500) {
                throw new Error('Internal Server Error');
            }
            return response.json();
        })
        .then(result => {
            if (result.status === 'Success') {
                document.getElementById('specialistProcess').style.display = 'block';
                alert('Successfully');
                cancelSpecialist()
            } else {
                alert("Something went wrong !!! " + result.error)
            }
        })
        .catch((error) => {
            console.error('Error:', error);
            if (error.message === 'Internal Server Error') {
                alert('An internal server error occurred. Please try again later.');
            } else {
                alert('An error occurred. Please try again.');
            }
        });

}

function exportSpecialist() {
    // Prevent default form submission
    this.event.preventDefault();


    fetch('/api/v1/specialist/export', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            if (response.status === 500) {
                throw new Error('Internal Server Error');
            }
            return response.json();
        })
        .then(result => {
            if (result.status === 'Success') {
                document.getElementById('specialistProcess').style.display = 'block';
                alert('Successfully');
                cancelSpecialist()
            } else {
                alert("Something went wrong !!! " + result.error)
            }
        })
        .catch((error) => {
            console.error('Error:', error);
            if (error.message === 'Internal Server Error') {
                alert('An internal server error occurred. Please try again later.');
            } else {
                alert('An error occurred. Please try again.');
            }
        });

}

function getDataFromInput() {
    // Get input values
    const RowNo = document.getElementById('RowNo').value;
    const CustomerNo = document.getElementById('CustomerNo').value;
    const CustomerNameTh = document.getElementById('CustomerNameTh').value;
    const ChkThaiName = document.getElementById('ChkThaiName').value;
    const ChkEngName = document.getElementById('ChkEngName').value;
    const CustomerNameEn = document.getElementById('CustomerNameEn').value;
    const Dob = document.getElementById('Dob').value;
    const CountryCode = document.getElementById('CountryCode').value;
    const CustomerType = document.getElementById('CustomerType').value;
    const Zipcode = document.getElementById('Zipcode').value;
    const Action = document.getElementById('Action').value;
    const OldAct = document.getElementById('OldAct').value;
    const Ovract = document.getElementById('Ovract').value;
    const Pob = document.getElementById('Pob').value;
    const ReasonCode = document.getElementById('ReasonCode').value;
    const RtnCustomer = document.getElementById('RtnCustomer').value;
    const SrcSeq = document.getElementById('SrcSeq').value;

    // Create request body
    return {
        RowNo: RowNo,
        CustomerNo: CustomerNo,
        CustomerNameTh: CustomerNameTh,
        ChkThaiName: ChkThaiName,
        ChkEngName: ChkEngName,
        CustomerNameEn: CustomerNameEn,
        Dob: Dob,
        CountryCode: CountryCode,
        CustomerType: CustomerType,
        Zipcode: Zipcode,
        Action: Action,
        OldAct: OldAct,
        Ovract: Ovract,
        Pob: Pob,
        ReasonCode: ReasonCode,
        RtnCustomer: RtnCustomer,
        SrcSeq: SrcSeq
    };
}

