import { createI18n } from 'vue-i18n'
import { GetCurrentLanguage } from '../../wailsjs/go/main/App'

const AppVersion = '2.0.5'

// 中文语言包
const zhCN = {
  // 应用信息
  app: {
    title: '剪存 • 剪贴板历史',
    name: '剪存',
    description: '剪贴板历史管理工具',
    version: AppVersion
  },

  // 登录页面
  login: {
    title: '请输入密码解锁',
    passwordPlaceholder: '请输入密码',
    unlockButton: '解锁',
    forgotPassword: '忘记密码？请删除数据库文件重置应用',
    dbLocation: '数据库位置: ~/.clipsave/clipboard.db',
    passwordRequired: '请输入密码',
    verifyFailed: '验证失败，请重试',
    unlockSuccess: '解锁成功！',
    passwordError: '密码错误，请重试',
    verifyError: '验证失败: {0}'
  },

  // 设置页面
  settings: {
    title: '设置',
    back: '返回',
    security: '安全设置',
    appPassword: '应用密码',
    passwordDesc: '设置密码后，每次打开应用需要输入密码',
    setPassword: '设置密码',
    changePassword: '修改密码',
    removePassword: '移除密码',
    removePasswordDesc: '移除密码后可直接打开应用',
    lock: '锁定',
    general: '通用设置',
    autoClean: '自动清理历史',
    autoCleanDesc: '自动删除超过指定天数的剪贴板历史',
    retentionDays: '保留天数',
    retentionDaysDesc: '历史记录保留的天数',
    hotkey: '全局快捷键',
    hotkeyDesc: '按下快捷键唤起应用窗口: {0}',
    recording: '录制中...',
    record: '录制',
    recordingPlaceholder: '请按下快捷键组合...',
    recordPlaceholder: '点击录制快捷键',
    clearAll: '全部清除',
    clearAllDesc: '清除所有剪贴板历史记录，此操作不可恢复',
    clearAllButton: '清除全部',
    interface: '界面设置',
    pageSize: '每页显示数量',
    pageSizeDesc: '列表中每次加载的记录数量',
    about: '关于',
    appName: '应用名称：',
    version: '版本号：',
    description: '描述：',
    language: '语言设置',
    languageDesc: '选择应用界面语言',
    backgroundMode: '后台运行',
    backgroundModeDesc: '开启后应用将在后台运行，不显示 Dock 图标'
  },

  // 密码设置对话框
  passwordDialog: {
    title: '设置应用密码',
    newPassword: '新密码',
    newPlaceholder: '请输入新密码',
    confirmPassword: '确认密码',
    confirmPlaceholder: '请再次输入密码',
    cancel: '取消',
    confirm: '确定',
    passwordRequired: '请输入密码',
    passwordMismatch: '两次输入的密码不一致',
    passwordTooShort: '密码长度至少4位',
    success: '密码设置成功！下次启动应用需要输入密码',
    error: '设置密码失败'
  },

  // 主界面
  main: {
    searchPlaceholder: '输入内容过滤...',
    filterAll: '所有类型',
    filterText: '文本',
    filterImage: '图片',
    filterFile: '文件',
    filterUrl: 'URL',
    filterColor: '颜色',
    filterJSON: 'JSON',
    listTitle: '列表',
    favorite: '收藏',
    unfavorite: '取消收藏',
    loading: '加载中...',
    emptyState: '暂无剪贴板历史',
    welcome: '欢迎使用 剪存！复制任何内容后，它将自动出现在这里。',
    source: '来源:',
    contentType: '内容类型:',
    charCount: '字符数:',
    wordCount: '单词数:',
    fileCount: '文件数:',
    createTime: '创建时间:',
    copy: '复制',
    delete: '删除',
    clipboardHistory: '剪贴板历史'
  },

  // 组件相关文本
  components: {
    // 文本组件
    text: {
      decodeUri: '解码 URI',
      decodeUnicode: '解码 Unicode',
      decodedText: '解码后文本',
      decodeFailed: '解码失败：{0}',
      translate: '翻译',
      translatedText: '翻译'
    },
    // 语言列表
    language: {
      zh: '中文',
      en: '英文',
      fr: '法语',
      de: '德语',
      es: '西班牙语',
      it: '意大利语',
      ru: '俄语',
      pt: '葡萄牙语',
      vi: '越南语',
      th: '泰语',
      ms: '马来语'
    },
    // 文件组件
    file: {
      fileNotExists: '（文件不存在）',
      openInFinder: '在 Finder 中打开'
    },
    // 图片组件
    image: {
      clipboardImage: '剪贴板图片',
      qrContent: '二维码内容：',
      copy: '复制',
      saveToLocal: '保存到本地',
      recognizing: '识别中...',
      recognizeQR: '识别二维码',
      qrGenerated: '二维码生成成功',
      qrGenerateFailed: '二维码生成失败',
      qrSaved: '二维码已保存',
      qrSaveFailed: '保存失败',
      qrCopied: '二维码已复制到剪贴板',
      qrCopyFailed: '复制失败'
    },
    // URL组件
    url: {
      openInBrowser: '在浏览器中打开',
      generating: '生成中...',
      generateQR: '生成二维码',
      generatedQR: '生成的二维码',
      saveQR: '保存二维码',
      copyQR: '复制二维码',
      qrGenerated: '二维码生成成功',
      qrGenerateFailed: '二维码生成失败',
      qrSaved: '二维码已保存',
      qrSaveFailed: '保存失败',
      qrCopied: '二维码已复制到剪贴板',
      qrCopyFailed: '复制失败',
      urlParams: 'URL 参数',
      key: '键',
      value: '值'
    },
    // 颜色组件
    color: {
      clickToSelect: '点击选择颜色',
      originalValue: '原始值:',
      rgb: 'RGB:',
      hex: 'HEX:',
      alpha: '透明度:',
      copied: '已复制'
    }
  },

  // 消息提示
  message: {
    copySuccess: '已复制到剪贴板',
    copyError: '复制失败: {0}',
    deleteConfirm: '确定要删除这条记录吗？',
    deleteConfirmTitle: '提示',
    deleteConfirmBtn: '确定',
    deleteCancelBtn: '取消',
    deleteSuccess: '删除成功',
    deleteError: '删除失败: {0}',
    openFinderSuccess: '已在 Finder 中打开文件',
    openFinderError: '打开文件失败: {0}',
    openUrlSuccess: '已在浏览器中打开链接',
    openUrlError: '打开链接失败: {0}',
    settingsSaved: '设置已保存',
    settingsError: '保存设置失败',
    hotkeyUpdated: '快捷键已更新',
    hotkeyError: '快捷键更新失败，请重试',
    clearConfirm: '确定要清除所有剪贴板历史记录吗？此操作不可恢复！',
    clearConfirmTitle: '确认清除',
    clearConfirmBtn: '确定清除',
    clearCancelBtn: '取消',
    clearProcessing: '正在清除所有记录...',
    clearSuccess: '已成功清除所有记录！',
    clearError: '清除失败: {0}',
    removePasswordConfirm: '移除密码后，将不再需要密码即可打开应用。确定要移除密码吗？',
    removePasswordTitle: '确认移除',
    removePasswordSuccess: '密码已移除',
    favoriteAdded: '已收藏',
    favoriteRemoved: '已取消收藏',
    favoriteError: '收藏操作失败'
  }
}

// 英文语言包
const enUS = {
  // 应用信息
  app: {
    title: 'ClipSave • Clipboard History',
    name: 'ClipSave',
    description: 'Clipboard History Management Tool',
    version: AppVersion
  },

  // 登录页面
  login: {
    title: 'Please enter password to unlock',
    passwordPlaceholder: 'Please enter password',
    unlockButton: 'Unlock',
    forgotPassword: 'Forgot password? Please delete database file to reset app',
    dbLocation: 'Database location: ~/.clipsave/clipboard.db',
    passwordRequired: 'Please enter password',
    verifyFailed: 'Verification failed, please try again',
    unlockSuccess: 'Unlocked successfully!',
    passwordError: 'Incorrect password, please try again',
    verifyError: 'Verification failed: {0}'
  },

  // 设置页面
  settings: {
    title: 'Settings',
    back: 'Back',
    security: 'Security Settings',
    appPassword: 'App Password',
    passwordDesc: 'After setting password, you need to enter password every time you open the app',
    setPassword: 'Set Password',
    changePassword: 'Change Password',
    removePassword: 'Remove Password',
    removePasswordDesc: 'Remove password to open app directly',
    lock: 'Lock',
    general: 'General Settings',
    autoClean: 'Auto Clean History',
    autoCleanDesc: 'Automatically delete clipboard history older than specified days',
    retentionDays: 'Retention Days',
    retentionDaysDesc: 'Number of days to keep history records',
    hotkey: 'Global Hotkey',
    hotkeyDesc: 'Press hotkey to bring up app window: {0}',
    recording: 'Recording...',
    record: 'Record',
    recordingPlaceholder: 'Please press hotkey combination...',
    recordPlaceholder: 'Click to record hotkey',
    clearAll: 'Clear All',
    clearAllDesc: 'Clear all clipboard history records, this operation cannot be undone',
    clearAllButton: 'Clear All',
    interface: 'Interface Settings',
    pageSize: 'Items Per Page',
    pageSizeDesc: 'Number of records to load each time in the list',
    about: 'About',
    appName: 'App Name:',
    version: 'Version:',
    description: 'Description:',
    language: 'Language Settings',
    languageDesc: 'Select application interface language',
    backgroundMode: 'Background Mode',
    backgroundModeDesc: 'When enabled, the app will run in the background without showing Dock icon'
  },

  // 密码设置对话框
  passwordDialog: {
    title: 'Set App Password',
    newPassword: 'New Password',
    newPlaceholder: 'Please enter new password',
    confirmPassword: 'Confirm Password',
    confirmPlaceholder: 'Please enter password again',
    cancel: 'Cancel',
    confirm: 'Confirm',
    passwordRequired: 'Please enter password',
    passwordMismatch: 'Passwords do not match',
    passwordTooShort: 'Password must be at least 4 characters',
    success: 'Password set successfully! Next time you start the app, you\'ll need to enter the password',
    error: 'Failed to set password'
  },

  // 主界面
  main: {
    searchPlaceholder: 'Enter content to filter...',
    filterAll: 'All Types',
    filterText: 'Text',
    filterImage: 'Image',
    filterFile: 'File',
    filterUrl: 'URL',
    filterColor: 'Color',
    filterJSON: 'JSON',
    listTitle: 'List',
    favorite: 'Favorites',
    unfavorite: 'Unfavorite',
    loading: 'Loading...',
    emptyState: 'No clipboard history',
    welcome: 'Welcome to ClipSave! After copying any content, it will automatically appear here.',
    source: 'Source:',
    contentType: 'Content Type:',
    charCount: 'Characters:',
    wordCount: 'Words:',
    fileCount: 'Files:',
    createTime: 'Created:',
    copy: 'Copy',
    delete: 'Delete',
    clipboardHistory: 'Clipboard History'
  },

  // 组件相关文本
  components: {
    // 文本组件
    text: {
      decodeUri: 'Decode URI',
      decodeUnicode: 'Decode Unicode',
      decodedText: 'Decoded Text',
      decodeFailed: 'Decode failed: {0}',
      translate: 'Translate',
      translatedText: 'Translation'
    },
    // 语言列表
    language: {
      zh: 'Chinese',
      en: 'English',
      fr: 'French',
      de: 'German',
      es: 'Spanish',
      it: 'Italian',
      ru: 'Russian',
      pt: 'Portuguese',
      vi: 'Vietnamese',
      th: 'Thai',
      ms: 'Malay'
    },
    // 文件组件
    file: {
      fileNotExists: '(File not exists)',
      openInFinder: 'Open in Finder'
    },
    // 图片组件
    image: {
      clipboardImage: 'Clipboard Image',
      qrContent: 'QR Code Content:',
      copy: 'Copy',
      saveToLocal: 'Save to Local',
      recognizing: 'Recognizing...',
      recognizeQR: 'Recognize QR Code',
      qrGenerated: 'QR code generated successfully',
      qrGenerateFailed: 'Failed to generate QR code',
      qrSaved: 'QR code saved',
      qrSaveFailed: 'Save failed',
      qrCopied: 'QR code copied to clipboard',
      qrCopyFailed: 'Copy failed'
    },
    // URL组件
    url: {
      openInBrowser: 'Open in Browser',
      generating: 'Generating...',
      generateQR: 'Generate QR Code',
      generatedQR: 'Generated QR Code',
      saveQR: 'Save QR Code',
      copyQR: 'Copy QR Code',
      qrGenerated: 'QR code generated successfully',
      qrGenerateFailed: 'Failed to generate QR code',
      qrSaved: 'QR code saved',
      qrSaveFailed: 'Save failed',
      qrCopied: 'QR code copied to clipboard',
      qrCopyFailed: 'Copy failed',
      urlParams: 'URL Parameters',
      key: 'Key',
      value: 'Value'
    },
    // 颜色组件
    color: {
      clickToSelect: 'Click to select color',
      originalValue: 'Original Value:',
      rgb: 'RGB:',
      hex: 'HEX:',
      alpha: 'Alpha:',
      copied: 'Copied'
    }
  },

  // 消息提示
  message: {
    copySuccess: 'Copied to clipboard',
    copyError: 'Copy failed: {0}',
    deleteConfirm: 'Are you sure you want to delete this record?',
    deleteConfirmTitle: 'Confirm',
    deleteConfirmBtn: 'Confirm',
    deleteCancelBtn: 'Cancel',
    deleteSuccess: 'Deleted successfully',
    deleteError: 'Delete failed: {0}',
    openFinderSuccess: 'File opened in Finder',
    openFinderError: 'Failed to open file: {0}',
    openUrlSuccess: 'Link opened in browser',
    openUrlError: 'Failed to open link: {0}',
    settingsSaved: 'Settings saved',
    settingsError: 'Failed to save settings',
    hotkeyUpdated: 'Hotkey updated',
    hotkeyError: 'Failed to update hotkey, please try again',
    clearConfirm: 'Are you sure you want to clear all clipboard history records? This operation cannot be undone!',
    clearConfirmTitle: 'Confirm Clear',
    clearConfirmBtn: 'Confirm Clear',
    clearCancelBtn: 'Cancel',
    clearProcessing: 'Clearing all records...',
    clearSuccess: 'All records cleared successfully!',
    clearError: 'Clear failed: {0}',
    removePasswordConfirm: 'After removing password, you won\'t need password to open the app. Are you sure you want to remove password?',
    removePasswordTitle: 'Confirm Remove',
    removePasswordSuccess: 'Password removed',
    favoriteAdded: 'Added to favorites',
    favoriteRemoved: 'Removed from favorites',
    favoriteError: 'Favorite action failed'
  }
}

// 法文语言包
const frFR = {
  // 应用信息
  app: {
    title: 'ClipSave • Historique du Presse-papiers',
    name: 'ClipSave',
    description: 'Outil de Gestion de l\'Historique du Presse-papiers',
    version: AppVersion
  },

  // 登录页面
  login: {
    title: 'Veuillez entrer le mot de passe pour déverrouiller',
    passwordPlaceholder: 'Veuillez entrer le mot de passe',
    unlockButton: 'Déverrouiller',
    forgotPassword: 'Mot de passe oublié ? Supprimez le fichier de base de données pour réinitialiser l\'application',
    dbLocation: 'Emplacement de la base de données : ~/.clipsave/clipboard.db',
    passwordRequired: 'Veuillez entrer le mot de passe',
    verifyFailed: 'Échec de la vérification, veuillez réessayer',
    unlockSuccess: 'Déverrouillé avec succès !',
    passwordError: 'Mot de passe incorrect, veuillez réessayer',
    verifyError: 'Échec de la vérification : {0}'
  },

  // 设置页面
  settings: {
    title: 'Paramètres',
    back: 'Retour',
    security: 'Paramètres de Sécurité',
    appPassword: 'Mot de Passe de l\'Application',
    passwordDesc: 'Après avoir défini un mot de passe, vous devrez l\'entrer à chaque ouverture de l\'application',
    setPassword: 'Définir le Mot de Passe',
    changePassword: 'Modifier le Mot de Passe',
    removePassword: 'Supprimer le Mot de Passe',
    removePasswordDesc: 'Supprimer le mot de passe pour ouvrir l\'application directement',
    lock: 'Verrouiller',
    general: 'Paramètres Généraux',
    autoClean: 'Nettoyage Automatique de l\'Historique',
    autoCleanDesc: 'Supprimer automatiquement l\'historique du presse-papiers plus ancien que le nombre de jours spécifié',
    retentionDays: 'Jours de Rétention',
    retentionDaysDesc: 'Nombre de jours pour conserver les enregistrements d\'historique',
    hotkey: 'Raccourci Global',
    hotkeyDesc: 'Appuyez sur le raccourci pour ouvrir la fenêtre de l\'application : {0}',
    recording: 'Enregistrement...',
    record: 'Enregistrer',
    recordingPlaceholder: 'Veuillez appuyer sur la combinaison de touches...',
    recordPlaceholder: 'Cliquez pour enregistrer le raccourci',
    clearAll: 'Tout Effacer',
    clearAllDesc: 'Effacer tous les enregistrements d\'historique du presse-papiers, cette opération est irréversible',
    clearAllButton: 'Tout Effacer',
    interface: 'Paramètres de l\'Interface',
    pageSize: 'Éléments par Page',
    pageSizeDesc: 'Nombre d\'enregistrements à charger à la fois dans la liste',
    about: 'À Propos',
    appName: 'Nom de l\'Application :',
    version: 'Version :',
    description: 'Description :',
    language: 'Paramètres de Langue',
    languageDesc: 'Sélectionner la langue de l\'interface de l\'application',
    backgroundMode: 'Mode Arrière-plan',
    backgroundModeDesc: 'Lorsqu\'il est activé, l\'application fonctionnera en arrière-plan sans afficher l\'icône du Dock'
  },

  // 密码设置对话框
  passwordDialog: {
    title: 'Définir le Mot de Passe de l\'Application',
    newPassword: 'Nouveau Mot de Passe',
    newPlaceholder: 'Veuillez entrer le nouveau mot de passe',
    confirmPassword: 'Confirmer le Mot de Passe',
    confirmPlaceholder: 'Veuillez entrer le mot de passe à nouveau',
    cancel: 'Annuler',
    confirm: 'Confirmer',
    passwordRequired: 'Veuillez entrer le mot de passe',
    passwordMismatch: 'Les mots de passe ne correspondent pas',
    passwordTooShort: 'Le mot de passe doit contenir au moins 4 caractères',
    success: 'Mot de passe défini avec succès ! La prochaine fois que vous démarrerez l\'application, vous devrez entrer le mot de passe',
    error: 'Échec de la définition du mot de passe'
  },

  // 主界面
  main: {
    searchPlaceholder: 'Entrer le contenu à filtrer...',
    filterAll: 'Tous les Types',
    filterText: 'Texte',
    filterImage: 'Image',
    filterFile: 'Fichier',
    filterUrl: 'URL',
    filterColor: 'Couleur',
    filterJSON: 'JSON',
    listTitle: 'Liste',
    favorite: 'Favoris',
    unfavorite: 'Retirer des favoris',
    loading: 'Chargement...',
    emptyState: 'Aucun historique du presse-papiers',
    welcome: 'Bienvenue dans ClipSave ! Après avoir copié du contenu, il apparaîtra automatiquement ici.',
    source: 'Source :',
    contentType: 'Type de Contenu :',
    charCount: 'Caractères :',
    wordCount: 'Mots :',
    fileCount: 'Fichiers :',
    createTime: 'Créé :',
    copy: 'Copier',
    delete: 'Supprimer',
    clipboardHistory: 'Historique du Presse-papiers'
  },

  // 组件相关文本
  components: {
    // 文本组件
    text: {
      decodeUri: 'Décoder URI',
      decodeUnicode: 'Décoder Unicode',
      decodedText: 'Texte Décodé',
      decodeFailed: 'Échec du décodage : {0}',
      translate: 'Traduire',
      translatedText: 'Traduction'
    },
    // 语言列表
    language: {
      zh: 'Chinois',
      en: 'Anglais',
      fr: 'Français',
      de: 'Allemand',
      es: 'Espagnol',
      it: 'Italien',
      ru: 'Russe',
      pt: 'Portugais',
      vi: 'Vietnamien',
      th: 'Thaï',
      ms: 'Malais'
    },
    // 文件组件
    file: {
      fileNotExists: '(Fichier n\'existe pas)',
      openInFinder: 'Ouvrir dans le Finder'
    },
    // 图片组件
    image: {
      clipboardImage: 'Image du Presse-papiers',
      qrContent: 'Contenu du Code QR :',
      copy: 'Copier',
      saveToLocal: 'Enregistrer Localement',
      recognizing: 'Reconnaissance...',
      recognizeQR: 'Reconnaître le Code QR',
      qrGenerated: 'Code QR généré avec succès',
      qrGenerateFailed: 'Échec de la génération du code QR',
      qrSaved: 'Code QR enregistré',
      qrSaveFailed: 'Échec de l\'enregistrement',
      qrCopied: 'Code QR copié dans le presse-papiers',
      qrCopyFailed: 'Échec de la copie'
    },
    // URL组件
    url: {
      openInBrowser: 'Ouvrir dans le Navigateur',
      generating: 'Génération...',
      generateQR: 'Générer le Code QR',
      generatedQR: 'Code QR Généré',
      saveQR: 'Enregistrer le Code QR',
      copyQR: 'Copier le Code QR',
      qrGenerated: 'Code QR généré avec succès',
      qrGenerateFailed: 'Échec de la génération du code QR',
      qrSaved: 'Code QR enregistré',
      qrSaveFailed: 'Échec de l\'enregistrement',
      qrCopied: 'Code QR copié dans le presse-papiers',
      qrCopyFailed: 'Échec de la copie',
      urlParams: 'Paramètres URL',
      key: 'Clé',
      value: 'Valeur'
    },
    // 颜色组件
    color: {
      clickToSelect: 'Cliquer pour sélectionner la couleur',
      originalValue: 'Valeur Originale :',
      rgb: 'RGB :',
      hex: 'HEX :',
      alpha: 'Alpha :',
      copied: 'Copié'
    }
  },

  // 消息提示
  message: {
    copySuccess: 'Copié dans le presse-papiers',
    copyError: 'Échec de la copie : {0}',
    deleteConfirm: 'Êtes-vous sûr de vouloir supprimer cet enregistrement ?',
    deleteConfirmTitle: 'Confirmer',
    deleteConfirmBtn: 'Confirmer',
    deleteCancelBtn: 'Annuler',
    deleteSuccess: 'Supprimé avec succès',
    deleteError: 'Échec de la suppression : {0}',
    openFinderSuccess: 'Fichier ouvert dans le Finder',
    openFinderError: 'Échec de l\'ouverture du fichier : {0}',
    openUrlSuccess: 'Lien ouvert dans le navigateur',
    openUrlError: 'Échec de l\'ouverture du lien : {0}',
    settingsSaved: 'Paramètres enregistrés',
    settingsError: 'Échec de l\'enregistrement des paramètres',
    hotkeyUpdated: 'Raccourci mis à jour',
    hotkeyError: 'Échec de la mise à jour du raccourci, veuillez réessayer',
    clearConfirm: 'Êtes-vous sûr de vouloir effacer tous les enregistrements d\'historique du presse-papiers ? Cette opération est irréversible !',
    clearConfirmTitle: 'Confirmer l\'Effacement',
    clearConfirmBtn: 'Confirmer l\'Effacement',
    clearCancelBtn: 'Annuler',
    clearProcessing: 'Effacement de tous les enregistrements...',
    clearSuccess: 'Tous les enregistrements effacés avec succès !',
    clearError: 'Échec de l\'effacement : {0}',
    removePasswordConfirm: 'Après avoir supprimé le mot de passe, vous n\'aurez plus besoin de mot de passe pour ouvrir l\'application. Êtes-vous sûr de vouloir supprimer le mot de passe ?',
    removePasswordTitle: 'Confirmer la Suppression',
    removePasswordSuccess: 'Mot de passe supprimé',
    favoriteAdded: 'Ajouté aux favoris',
    favoriteRemoved: 'Retiré des favoris',
    favoriteError: 'Échec de l’action des favoris'
  }
}

// 阿拉伯语语言包
const arSA = {
  // 应用信息
  app: {
    title: 'ClipSave • سجل الحافظة',
    name: 'ClipSave',
    description: 'أداة إدارة سجل الحافظة',
    version: AppVersion
  },

  // 登录页面
  login: {
    title: 'يرجى إدخال كلمة المرور لإلغاء القفل',
    passwordPlaceholder: 'يرجى إدخال كلمة المرور',
    unlockButton: 'إلغاء القفل',
    forgotPassword: 'نسيت كلمة المرور؟ احذف ملف قاعدة البيانات لإعادة تعيين التطبيق',
    dbLocation: 'موقع قاعدة البيانات: ~/.clipsave/clipboard.db',
    passwordRequired: 'يرجى إدخال كلمة المرور',
    verifyFailed: 'فشل التحقق، يرجى المحاولة مرة أخرى',
    unlockSuccess: 'تم إلغاء القفل بنجاح!',
    passwordError: 'كلمة المرور غير صحيحة، يرجى المحاولة مرة أخرى',
    verifyError: 'فشل التحقق: {0}'
  },

  // 设置页面
  settings: {
    title: 'الإعدادات',
    back: 'رجوع',
    security: 'إعدادات الأمان',
    appPassword: 'كلمة مرور التطبيق',
    passwordDesc: 'بعد تعيين كلمة المرور، ستحتاج إلى إدخالها في كل مرة تفتح فيها التطبيق',
    setPassword: 'تعيين كلمة المرور',
    changePassword: 'تغيير كلمة المرور',
    removePassword: 'إزالة كلمة المرور',
    removePasswordDesc: 'إزالة كلمة المرور لفتح التطبيق مباشرة',
    lock: 'قفل',
    general: 'الإعدادات العامة',
    autoClean: 'تنظيف السجل التلقائي',
    autoCleanDesc: 'حذف سجل الحافظة تلقائياً الأقدم من الأيام المحددة',
    retentionDays: 'أيام الاحتفاظ',
    retentionDaysDesc: 'عدد الأيام للاحتفاظ بسجلات السجل',
    hotkey: 'اختصار عام',
    hotkeyDesc: 'اضغط على الاختصار لفتح نافذة التطبيق: {0}',
    recording: 'جاري التسجيل...',
    record: 'تسجيل',
    recordingPlaceholder: 'يرجى الضغط على مجموعة المفاتيح...',
    recordPlaceholder: 'انقر لتسجيل الاختصار',
    clearAll: 'مسح الكل',
    clearAllDesc: 'مسح جميع سجلات سجل الحافظة، هذه العملية لا يمكن التراجع عنها',
    clearAllButton: 'مسح الكل',
    interface: 'إعدادات الواجهة',
    pageSize: 'العناصر في الصفحة',
    pageSizeDesc: 'عدد السجلات المراد تحميلها في كل مرة في القائمة',
    about: 'حول',
    appName: 'اسم التطبيق:',
    version: 'الإصدار:',
    description: 'الوصف:',
    language: 'إعدادات اللغة',
    languageDesc: 'اختر لغة واجهة التطبيق',
    backgroundMode: 'وضع الخلفية',
    backgroundModeDesc: 'عند التمكين، سيعمل التطبيق في الخلفية دون عرض أيقونة Dock'
  },

  // 密码设置对话框
  passwordDialog: {
    title: 'تعيين كلمة مرور التطبيق',
    newPassword: 'كلمة المرور الجديدة',
    newPlaceholder: 'يرجى إدخال كلمة المرور الجديدة',
    confirmPassword: 'تأكيد كلمة المرور',
    confirmPlaceholder: 'يرجى إدخال كلمة المرور مرة أخرى',
    cancel: 'إلغاء',
    confirm: 'تأكيد',
    passwordRequired: 'يرجى إدخال كلمة المرور',
    passwordMismatch: 'كلمات المرور غير متطابقة',
    passwordTooShort: 'يجب أن تكون كلمة المرور 4 أحرف على الأقل',
    success: 'تم تعيين كلمة المرور بنجاح! في المرة القادمة التي تفتح فيها التطبيق ستحتاج إلى إدخال كلمة المرور',
    error: 'فشل في تعيين كلمة المرور'
  },

  // 主界面
  main: {
    searchPlaceholder: 'أدخل المحتوى للتصفية...',
    filterAll: 'جميع الأنواع',
    filterText: 'نص',
    filterImage: 'صورة',
    filterFile: 'ملف',
    filterUrl: 'رابط',
    filterColor: 'لون',
    filterJSON: 'JSON',
    listTitle: 'قائمة',
    favorite: 'المفضلة',
    unfavorite: 'إزالة من المفضلة',
    loading: 'جاري التحميل...',
    emptyState: 'لا يوجد سجل حافظة',
    welcome: 'مرحباً بك في ClipSave! بعد نسخ أي محتوى، سيظهر تلقائياً هنا.',
    source: 'المصدر:',
    contentType: 'نوع المحتوى:',
    charCount: 'الأحرف:',
    wordCount: 'الكلمات:',
    fileCount: 'الملفات:',
    createTime: 'تاريخ الإنشاء:',
    copy: 'نسخ',
    delete: 'حذف',
    clipboardHistory: 'سجل الحافظة'
  },

  // 组件相关文本
  components: {
    // 文本组件
    text: {
      decodeUri: 'فك تشفير URI',
      decodeUnicode: 'فك تشفير Unicode',
      decodedText: 'النص المفكوك',
      decodeFailed: 'فشل فك التشفير: {0}',
      translate: 'ترجمة',
      translatedText: 'الترجمة'
    },
    // 语言列表
    language: {
      zh: 'الصينية',
      en: 'الإنجليزية',
      fr: 'الفرنسية',
      de: 'الألمانية',
      es: 'الإسبانية',
      it: 'الإيطالية',
      ru: 'الروسية',
      pt: 'البرتغالية',
      vi: 'الفيتنامية',
      th: 'التايلاندية',
      ms: 'الماليزية'
    },
    // 文件组件
    file: {
      fileNotExists: '(الملف غير موجود)',
      openInFinder: 'فتح في Finder'
    },
    // 图片组件
    image: {
      clipboardImage: 'صورة الحافظة',
      qrContent: 'محتوى رمز الاستجابة السريعة:',
      copy: 'نسخ',
      saveToLocal: 'حفظ محلياً',
      recognizing: 'جاري التعرف...',
      recognizeQR: 'التعرف على رمز الاستجابة السريعة',
      qrGenerated: 'تم إنشاء رمز الاستجابة السريعة بنجاح',
      qrGenerateFailed: 'فشل في إنشاء رمز الاستجابة السريعة',
      qrSaved: 'تم حفظ رمز الاستجابة السريعة',
      qrSaveFailed: 'فشل في الحفظ',
      qrCopied: 'تم نسخ رمز الاستجابة السريعة إلى الحافظة',
      qrCopyFailed: 'فشل في النسخ'
    },
    // URL组件
    url: {
      openInBrowser: 'فتح في المتصفح',
      generating: 'جاري الإنشاء...',
      generateQR: 'إنشاء رمز الاستجابة السريعة',
      generatedQR: 'رمز الاستجابة السريعة المنشأ',
      saveQR: 'حفظ رمز الاستجابة السريعة',
      copyQR: 'نسخ رمز الاستجابة السريعة',
      qrGenerated: 'تم إنشاء رمز الاستجابة السريعة بنجاح',
      qrGenerateFailed: 'فشل في إنشاء رمز الاستجابة السريعة',
      qrSaved: 'تم حفظ رمز الاستجابة السريعة',
      qrSaveFailed: 'فشل في الحفظ',
      qrCopied: 'تم نسخ رمز الاستجابة السريعة إلى الحافظة',
      qrCopyFailed: 'فشل في النسخ',
      urlParams: 'معاملات الرابط',
      key: 'المفتاح',
      value: 'القيمة'
    },
    // 颜色组件
    color: {
      clickToSelect: 'انقر لاختيار اللون',
      originalValue: 'القيمة الأصلية:',
      rgb: 'RGB:',
      hex: 'HEX:',
      alpha: 'الشفافية:',
      copied: 'تم النسخ'
    }
  },

  // 消息提示
  message: {
    copySuccess: 'تم النسخ إلى الحافظة',
    copyError: 'فشل في النسخ: {0}',
    deleteConfirm: 'هل أنت متأكد من حذف هذا السجل؟',
    deleteConfirmTitle: 'تأكيد',
    deleteConfirmBtn: 'تأكيد',
    deleteCancelBtn: 'إلغاء',
    deleteSuccess: 'تم الحذف بنجاح',
    deleteError: 'فشل في الحذف: {0}',
    openFinderSuccess: 'تم فتح الملف في Finder',
    openFinderError: 'فشل في فتح الملف: {0}',
    openUrlSuccess: 'تم فتح الرابط في المتصفح',
    openUrlError: 'فشل في فتح الرابط: {0}',
    settingsSaved: 'تم حفظ الإعدادات',
    settingsError: 'فشل في حفظ الإعدادات',
    hotkeyUpdated: 'تم تحديث الاختصار',
    hotkeyError: 'فشل في تحديث الاختصار، يرجى المحاولة مرة أخرى',
    clearConfirm: 'هل أنت متأكد من مسح جميع سجلات سجل الحافظة؟ هذه العملية لا يمكن التراجع عنها!',
    clearConfirmTitle: 'تأكيد المسح',
    clearConfirmBtn: 'تأكيد المسح',
    clearCancelBtn: 'إلغاء',
    clearProcessing: 'جاري مسح جميع السجلات...',
    clearSuccess: 'تم مسح جميع السجلات بنجاح!',
    clearError: 'فشل في المسح: {0}',
    removePasswordConfirm: 'بعد إزالة كلمة المرور، لن تحتاج إلى كلمة مرور لفتح التطبيق. هل أنت متأكد من إزالة كلمة المرور؟',
    removePasswordTitle: 'تأكيد الإزالة',
    removePasswordSuccess: 'تم إزالة كلمة المرور',
    favoriteAdded: 'تمت الإضافة إلى المفضلة',
    favoriteRemoved: 'تمت الإزالة من المفضلة',
    favoriteError: 'فشل إجراء المفضلة'
  }
}

// 创建i18n实例
const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN', // 默认语言
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
    'fr-FR': frFR,
    'ar-SA': arSA
  }
})

// 从后端获取当前语言并设置
async function initLanguage() {
  try {
    const currentLang = await GetCurrentLanguage()
    i18n.global.locale.value = currentLang as any
  } catch (error) {
    console.error('Failed to get current language:', error)
  }
}

// 初始化语言
initLanguage()

export default i18n
