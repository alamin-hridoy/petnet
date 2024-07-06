function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

function IsImageFromUrl(url){
    var imgs = ["image/jpeg", "image/png", "image/gif"];
    var http = new XMLHttpRequest();
    http.open("HEAD", url, false);
    http.send();
    var type = false;
    if(http.readyState==4 && http.status==200) {
        var h = http.getAllResponseHeaders();
        var t = h.slice(h.indexOf("content-type:")+14,);
        if(t.includes(';')){
            var x = t.split(';');
            t = x[0];
        }
        var s = t.split('\n');
        type = s[0].replace(/[\r\n]/g, "");
        return imgs.includes(type);
    }
    return type;
}

function isFileImage(file) {
    return new Promise((resolve) => {
        const image = new Image();

        image.onload = function() {
            resolve(!!this.width);
        };

        window.setTimeout(()=>resolve(false), 5*100)

        image.src = URL.createObjectURL(file);
    })
}

class FileUpload {
    constructor(
        partname,
        maxfsize,
        endings,
        endingsString,
        uploadDesc,
        uploadedFiles,
        errors,
        options,
    ) {
        if (options == null) {
            options = {};
        }
        this.doPreview = !!options.previewImages
        this.multiple = !!options.multiple
        this.partname = partname
        this.index = options.index
        this.uploadCallback = options.uploadCallback

        this.maxfsize = maxfsize
        this.endings = endings
        this.endingsString = endingsString

        this.uploadDesc = uploadDesc
        this.uploadedFiles = uploadedFiles
        this.errors = errors
        this.hasErrors = false

        this.inputCnt = 1;

        this.mainElementInputID = this.generateElementID('file-upload')
        this.elementLabelID = this.generateElementID('upload-file-label')
        this.elementErrorsID = this.generateElementID('errors')
        this.elementFilelistID = this.generateElementID('filelist')
        this.elementUploadedFileListID = this.generateElementID('filelist-uploaded')
        this.elementUploadedFiletypesID = this.generateElementID('upload-filetypes') // used from outside
        this.elementInputID = (i) => this.generateElementID('input-'+i)
        this.elementRemoveFileID = (i) => this.generateElementID('removefile-'+i)
        this.elementRemoveUploadedFileID = (i) => this.generateElementID('removefile-uploaded-'+i)

        this.elementUploadedInputName = `uploaded-${this.partname}`
        this.elementInputName = this.partname;
        
        (async ()=> {
            try {
                this.getMainElementInput().innerHTML = this.createHTML();

                this.createAttachInput(0)

                await this.renderFiles();
            } catch (error) {
                console.error(`cannot initialize ${this.partname}`, error)
            }
        })()
    }

    generateElementID(prefix) {
        return `${prefix}-${this.partname}${this.index ? '-' + this.index : ''}`
    }

    getMainElementInput() {
        return document.getElementById(this.mainElementInputID)
    }

    getElementLabel() {
        return document.getElementById(this.elementLabelID)
    }

    getElementInput(i) {
        return document.getElementById(this.elementInputID(i))
    }

    getElementErrors() {
        return document.getElementById(this.elementErrorsID)
    }

    getElementFileList() {
        return document.getElementById(this.elementFilelistID)
    }

    getElementUploadedFileList() {
        return document.getElementById(this.elementUploadedFileListID )
    }

    getElementsRemoveFile(i) {
        return document.querySelectorAll('.'+this.elementRemoveFileID(i))
    }

    getElementRemoveUploadedFile(i) {
        return document.getElementById(this.elementRemoveUploadedFileID(i))
    }

    async removeFile(i, fName) {
        if (this.multiple) {
            let fls = this.getElementInput(i).files;
            if (fls.length == 0 ){
                this.getElementInput(i).remove();
                this.getElementInput(i).value = "";
            }
            if (fls.length > 0) {
                const dt = new DataTransfer();
                for (let file of fls){
                    if (file.name != fName) {
                        dt.items.add(file);
                    }
                }
                this.getElementInput(i).files = dt.files;
            }
        } else {
            this.getElementInput(i).value = "";
        }
        await this.renderFiles();
    }

    async removeUploadedFile(i) {
        this.uploadedFiles[i]=null;
        await this.renderFiles();
    }

    checkFileName(fileName) {
        for (const ending of this.endings) {
            if (fileName.toLowerCase().endsWith("."+ending)) {
                return true
            }
        }
        return false
    }

    fileEndingsString() {
        return this.endingsString
    }

    getAllFiles() {
        let files = {}
        for (let i = 0; i < this.inputCnt; i++) {
            const inp = this.getElementInput(i)
            if (inp != null) {
                const fileList = Array.from(this.getElementInput(i).files)
                files[i] = fileList;
            }
        }
        return files;
    }

    async renderFiles() {
        this.renderUploadedFiles();
        await this.renderNewFiles();
    }

    renderUploadedFiles() {
        let html = '';

        for (let i = 0; i < this.uploadedFiles.length ; i++) {

            const f = this.uploadedFiles[i]
            if (f != null) {
                let img = `<img
                        src="/images/file.png"
                        alt="file"
                        width="18"
                        class="ml-4 mt-2 mb-2 align-middle inline-block flex-initial"
                    />`
                if (this.doPreview) {

                    if(f.video=="true")
                    {
                        img=`<video preload="metadata" class="ml-4 mt-2 mb-2 align-middle inline-block flex-initial max-h-48 max-w-48">
                        <source src="${f.url}#t=0.1" type="video/mp4">
                        <source src="${f.url}#t=0.1" type="video/mov">
                        Your browser does not support the video tag.
                        </video>`
                    } 
                    else {
                        let imgUrl = f.url;
                        if(!IsImageFromUrl(f.url)) {
                            imgUrl = "/images/file.png";
                        }
                        img = `<img
                        src="${imgUrl}"
                        alt="file"
                        class="ml-4 mt-2 mb-2 align-middle inline-block flex-initial max-h-8 max-w-8"
                    />`
                    }
                    
                }

                html += `<div class="mb-2 p-4 border rounded border-grey-light-5 flex items-center remove_file_wrap">
                ${img}
                <div class="inline-block align-middle ml-4 text-blue-dark-5 flex-initial">
                    <a class="text-blue cursor-pointer"
                        target="_blank"
                        href="${f.url}">
                        ${f.name}
                    </a>
                </div>
                <div id="${this.elementRemoveUploadedFileID(i)}"
                    class="remove_file inline-block align-middle ml-4 text-right text-red-500 cursor-pointer hover:underline flex-1 mr-2" data-fileid="${f.input}" data-file="${f.name}"
                >Remove
                </div>
                <input name="${this.elementUploadedInputName}" type="hidden"
                    value="${f.input}">
                </div>`
                
            }
        }

        this.getElementUploadedFileList().innerHTML = html;



        for (let i = 0; i < this.uploadedFiles.length; i++) {
            const f = this.uploadedFiles[i]
            if (f != null) {
                const l = this.getElementRemoveUploadedFile(i);
                l.addEventListener('click', () => this.removeUploadedFile(i));
            }
        }
    }

    async renderNewFiles() {
        const docErrors = this.getElementErrors();
        const files = this.getAllFiles();
        let html = '';
        this.hasErrors = false;
        // map for detecting duplicate file names
        let fileNames = new Map();
        let supportedFiles = this.endings.join(", ");
        for (const i in files) {
            const fArray = files[i];
            for (let j=0; j<fArray.length; j++) {
                const f = fArray[j];
                let name = f.name;
                const sizeB = f.size;
                const sizeFormatted = formatBytes(sizeB);
                let hasFileBefore = fileNames.has(f.name);

                if (name.length > 43) {
                    name = name.substring(0, 20) + "..." + name.substring(name.length - 20)
                }

                if (sizeB > this.maxfsize * 1024 * 1024) {
                    docErrors.innerText = `File is to large. Maximum file size is ${this.maxfsize}MB.`;
                    docErrors.classList.remove('hidden');
                    this.hasErrors = true;
                } else if (!this.checkFileName(name)) {
                    docErrors.innerText = `Invalid file type. Only ${supportedFiles} file types are allowed`;
                    docErrors.classList.remove('hidden');
                    this.hasErrors = true
                } else if (hasFileBefore) {
                    this.hasErrors = true;
                    docErrors.innerText = `The file has already been uploaded. Please select another file.`;
                    docErrors.classList.remove('hidden');
                } else {
                    // Valid files
                    if (!hasFileBefore) {
                        fileNames.set(f.name, true)
                    }

                    this.hasErrors = false;
                    docErrors.innerText = "";
                    docErrors.classList.add('hidden');
                }

                // Need to remove created input by function createAttachInput on file change event.
                // Else it keeps invalid file in input for multiple files.
                // TODO: file validations should be done on file change event before creating new input
                if (this.hasErrors) {
                    if (this.multiple) {
                        this.getElementInput(i).remove();
                    } else {
                        this.getElementInput(i).value = ""
                    }

                    continue
                }
                let elm = document.getElementsByClassName("all-errs-"+this.getElementInput(i).name);
                for (let i = 0; i < elm.length; i++) {
                    elm[i].innerHTML = "";
                }
                const fHtml = await this.previewFileHTML(i, j, f, name, sizeFormatted)
                html += fHtml
            }
        }

        this.getElementFileList().innerHTML = html;

        for (const i in files) {
            const l = this.getElementsRemoveFile(i);
            l.forEach(e => e.addEventListener('click', (ab) => {
                let prntElm = ab.target.parentElement.parentElement;
                let fName = prntElm.firstElementChild.firstElementChild.firstChild.data;
                prntElm.parentElement.remove();
                this.removeFile(i, fName.trim());
            }));
        }

        if (this.uploadCallback) {
            this.uploadCallback();
        }
    }

    uploadErrorClass() {
        return this.errors !== "" ? "border-red" : ""
    }

    errorHiddenClass() {
        return this.errors === "" ? "hidden" : ""
    }

    async previewFileHTML(i, j, file, name, sizeFormatted) {
        let img = `<img
                src="/images/file.png"
                alt="file"
                width="18"
                class="ml-4 mt-2 mb-2 align-middle inline-block flex-initial"
            />`
            
        if (this.doPreview) {
            const isImage = await isFileImage(file)
            if (isImage) {
                img = `<img
                    src="${URL.createObjectURL(file)}"
                    alt="file"
                    class="ml-4 mt-2 mb-2 max-h-8 max-w-8 align-middle inline-block flex-initial"
                />`
            }
        }
        if(!this.hasErrors){
            return `<div class="mb-2 p-4 border rounded border-grey-light-5 flex flex-wrap items-center">
            ${img}
            <div class="flex justify-between items-center">
                <div class="inline-block align-middle ml-4 text-blue-dark-5 flex-initial">
                    <span">${name}</span>
                    <span class="text-grey">${sizeFormatted}</span>
                </div>
                <div
                    class="${this.elementRemoveFileID(i)} inline-block align-middle ml-4 text-right text-red-500 cursor-pointer hover:underline flex-1 mr-2"
                >
                <span class="hidden sm:inline-block">Remove</span>
            </div>
            </div>
            <div class="${this.elementRemoveFileID(i)} cursor-pointer ml-3">
            <span class="sm:hidden inline-block">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7" viewBox="0 0 20 20" fill="#EF4444">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
                </span>
            </div>
        </div>`
        }
        return ''
    }

    createInput(i) {
        const fileInput = document.createElement('input');
        fileInput.setAttribute('name', this.elementInputName);
        fileInput.setAttribute('type', 'file');
        if (this.multiple) {
            fileInput.setAttribute('multiple', 'multiple')
        }
        fileInput.setAttribute('id', this.elementInputID(i))
        fileInput.setAttribute('accept', this.endings.map(e=>`.${e}`).join(','))
        fileInput.classList.add(...['w-0', 'opacity-0', 'hidden', 'upin']);
        return fileInput;
    }

    createAttachInput(i) {
        const inp = this.createInput(i);
        inp.addEventListener('change', async () => {
            if (this.multiple) {
                this.createAttachInput(this.inputCnt);
                this.inputCnt++;
            } else {
                this.uploadedFiles=[];
            }
            const docErrors = this.getElementErrors()
            docErrors.classList.add('hidden')
            await this.renderFiles()
        })
        this.getElementLabel().append(inp);
        this.getElementLabel().setAttribute('for', this.elementInputID(i));
    }

    createHTML() {
        
        let elements = `<div>
                    <div id="${this.elementUploadedFileListID}"></div>

                    <div id="${this.elementFilelistID}"></div>

                    <div class="flex items-center justify-between p-4 border border-gray-300 rounded bg-fileupload ${this.uploadErrorClass()}">
                        <div class="flex flex-col">
                            <div class="font-semibold text-sm text-gray-400">${this.uploadDesc}</div>
                    </div>
                
                    <div class="upload-file relative inline-block ml-2">
                        <label
                                id="${this.elementLabelID}"
                                class="inline-block py-1 px-8 cursor-pointer border rounded-md border-blue-dark-5 text-blue-dark-5 transition duration-300 ease-in-out hover:bg-blue-dark-5 font-bold text-petnetblue  border-petnetblue"
                        >
                            Upload
                        </label>
                    </div>
                </div>
                <p id="${this.elementErrorsID}" class="text-petnetpink bg-white pt-2  ${this.errorHiddenClass()}
                    style="white-space: pre-line">${this.errors}</p>
            </div>`;
            
        return elements;
    }
}
