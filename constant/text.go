package constant

const (
	CreateAccountSuccess     = "Congratulations\\! Your Meta wallet has been created \\. \n\nThe wallet address is\\: `%s`\\.\nPin Code is `%s`\\.\nYour Wallet Pin Code is the only way to access your crypto asset in MetaWallet and CAN NOT be recovered if lost\\."
	GetAccountSuccess        = "Your MetaWallet address is: `%s` \\."
	ButtonForwardPrivateChat = "FORWARD TO CONTINUE"
	Participate              = "Participate"
	SwitchPrivate            = "✳️ Please forward to private chat with bot for detail."
	NoStart                  = "Please use /start command first"
)

const (
	DiscordSavePinCode = `This message will disappear soon, please save it in time or use the /change_pin_code command to update the Pin Code

This message will also be privately chatted to you, and the private chat message will not disappear

If you don't receive a private chat message, your privacy settings may have prevented the bot from privately chatting with you. And you will not receive Pin Code private chat again.`

	AllianceStart = "ℹ️ User Guide\n\U0001F973🙌 Welcome\\! " +
		"Tristan MetaWallet is the embedded multi\\-chain account of Tristan Alliance — " +
		"a marketing and operating solution for Web3 projects\\. \n" +
		"With MetaWallet you can\\:\n" +
		"💸 Receive token rewards from community activities\n💵 " +
		"Transfer your assets to other third\\-party wallet\n" +
		"🚀 Get QTZ by taking part in activities to unlock potential utility after the Tristan alliance alpha test\\.\n \n" +
		"⚙️ Commands\n" +
		"/start Create your MetaWallet and get the user guide\\.\n" +
		"/change\\_pin\\_code Change pin code of your MetaWallet address\\.\n" +
		"/export\\_private Export the private key of your MetaWallet, and you can add this address to Metamask or other third\\-party wallet with it\\.\n" +
		"/replace\\_private Import a new account to replace the old one\\. Remember to backup your old account in advance\\, or you can't recover it once it is replaced\\.\n" +
		"/address Check your MetaWallet address\\.\n" +
		"/balance Get details of your MetaWallet balance for following assets\\: Crypto and NFTs\n" +
		"/transfer Transfer assets to certain address\\.\n" +
		"/add\\_token\\_balance Add specific token to display under \"/balance\" command\n \nYou always need your Pincode to execute your transfer transaction\\. *REMEMBER YOUR PINCODE AND KEEP IT SAFE*\\!\\!\\!\n"
)
