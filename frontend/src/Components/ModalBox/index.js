import { Modal, ModalContent } from '../ModalBox/Styles'


const ModalBox = ({open,onClose,children}) => {
    if (!open) return null

    return (
        <Modal onClick={onClose}>
            <ModalContent onClick={(e) => e.stopPropagation()}>
                {children}
            </ModalContent>
        </Modal>
    )
}

export default ModalBox

